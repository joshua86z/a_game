package controllers

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"encoding/binary"
	"fmt"
	"libs/log"
	"libs/ssdb"
	"libs/token"
	"strings"
	"time"
	//requestLog "game/models/logs"
	"encoding/json"
	"models"
	"protodata"
	"sync"
)

// 运行变量
var (
	online          int
	gameToken       *token.Token
	playLock        sync.Mutex
	playerMap       map[int64]*Connect
	request_log_map map[int32]string
)

func init() {
	gameToken = token.NewToken(tokenAdapter{})
	playerMap = make(map[int64]*Connect)
	request_log_map = make(map[int32]string)
	CountOnline()
}

// Client Connect
type Connect struct {
	Role    *models.RoleModel
	Conn    *websocket.Connect
	Chan    chan string
	Request *protodata.CommectRequest
}

func (this *Connect) Send(code protodata.StatusCode, value interface{}, err error) {
	if err != nil {
		value = fmt.Sprintf("%v", err)
	}
	this.send(ReturnStr(this.Request, code, value))
	this.Request = nil
}

func (this *Connect) send(s string) {
	this.Chan <- s
}

func (this *Connect) pushToClient() {
	go func() {
		for s := range this.Chan {
			var buf = make([]byte, 4)
			binary.LittleEndian.PutUint32(buf, uint32(len(s)))

			Buffer := bytes.NewBuffer(buf)
			Buffer.WriteString(s)

			if err := websocket.Message.Send(this.Conn, Buffer.Bytes()); err != nil {
				log.Warn("Can't send msg. %v", err)
			} else {
				log.Info("Send Success")
			}
		}
	}()
}

func (this *Connect) InMap(uid int64) {
	playLock.Lock()
	if _, ok := playerMap[uid]; !ok {
		playerMap[this.Role.Uid] = this
	}
	playLock.Unlock()
}

func (this *Connect) Close() {
	playLock.Lock()
	delete(playerMap, this.Role.Uid)
	playLock.Unlock()
	this.Conn.Close()
}

// 从客户端读取信息
func (this *Connect) PullFromClient() {

	for {
		// receive from ws Connect
		var content string
		if err := websocket.Message.Receive(this.Conn, &content); err != nil {
			if err.Error() == "EOF" {
				log.Info("Conn receive EOF")
			} else {
				log.Error("Can't receive message. %v", err)
			}
			return
		}

		// **************** 其它接口 **************** //
		if strings.HasPrefix(content, "20140709_allhero_") {
			request := strings.Replace(content, "20140709_allhero_", "", len(content))
			this.OtherRequest([]byte(request))
			return
		}
		// **************** 支付专用 **************** //

		beginTime := time.Now()
		log.Info(" Begin ")

		// parse proto message

		request, err := ParseContent(content)
		if err != nil {
			log.Error("Parse client request error. %v", err)
			this.send(ReturnStr(request, lineNum(), fmt.Sprintf("客户端错误:%v", err)))
			continue
		} else {
			this.Request = request
		}

		index := request.GetCmdId()

		// Panic recover
		defer func() {
			if err := recover(); err != nil {
				log.Critical("Panic occur. %v", err)
				this.Send(lineNum(), nil, err)
			}
		}()

		if index != 10000 {
			// Check Login status
			if request.GetTokenStr() == "" {
				this.Send(protodata.StatusCode_INVALID_TOKEN, nil, nil)
				continue
			}
			uid, _ := gameToken.GetUid(request.GetTokenStr())
			if uid == 0 {
				this.Send(protodata.StatusCode_INVALID_TOKEN, nil, nil)
				continue
			} else {
				if this.Role == nil {
					this.Role = models.NewRoleModel(uid)
				} else if this.Role.Uid != uid {
					this.Send(protodata.StatusCode_INVALID_TOKEN, nil, nil)
					continue
				}
			}
		}

		// 执行命令
		if function, ok := handlers[index]; ok {
			log.Info("Exec %v -> %s (uid:%d)", index, function, RoleModel.Uid)
			function()
		} else {
			this.Send(lineNum(), nil, fmt.Errorf("没有这方法 index : %d", index))
			continue
		}

		// 执行命令
		//handler()

		execTime := time.Now().Sub(beginTime)
		if execTime.Seconds() > 0.1 {
			//慢日志
			log.Warn("Slow Exec -> %s, time is %v second", handlerNames[index], execTime.Seconds())
		} else {
			log.Info("time is %v second", execTime.Seconds())
		}

		// Send response to client
		//this.Send(response)
		//this.Send <-

		if this.Role.Uid != 0 && index > 0 {
			//玩家操作记录
			if _, ok := request_log_map[index]; ok {
				//				requestLog.InsertLog(player.UniqueId, index)
			}
		}
	}
}

func (this *Connect) Function(index int32) func() {
	switch index {
	case 10103:
		return this.Login
	case 10104:
		return this.UserDataRequest
	case 10105:
		return this.SetRoleName
	case 10106:
		return this.RandomName
	case 10107:
		return this.BuyStaminaRequest
	case 10108:
		return this.BuyCoinRequest
	case 10110:
		return this.GeneralLevelUp
	case 10111:
		return this.Buyeneral
	case 10112:
		return this.MailList
	case 10113:
		return this.MailRewardRequest
	case 10114:
		return this.ItemLevelUp
	default:
		return func() {}
	}
}

func (this *Connect) OtherRequest(request []byte) {

	data := make(map[string]string)
	json.Unmarshal(request, &data)

	if data["cmd"] == "9000" {
		if OrderModel, err := models.NewOrderModel(data["orderId"]); err != nil {
			log.Error("Pay Error %v", err)
		} else {
			OrderModel.Confirm()
		}
	}
}

func Handler(ws *websocket.Connect) {

	online++
	// New Connectect
	Connect := &Connect{
		Conn: ws,
		Chan: make(chan string, 10),
	}

	Connect.pushToClient()
	Connect.PullFromClient()
	Connect.Close()

	online--
}

func SendMessage(uid int64, message string) error {
	playLock.Lock()
	if Connect, ok := playerMap[uid]; !ok {
		return fmt.Errorf("uid : %d not online", uid)
	} else {
		Connect.send(message)
	}
	playLock.Unlock()
	return nil
}

func CountOnline() {
	go func() {
		//	t := time.Tick(time.Second * 5)
		t := time.Tick(time.Minute * 5)
		for {
			select {
			case <-t:
				//	fmt.Println("online num : ", online)
				models.DB().Exec("INSERT INTO `stat_online`(`online_num`,`online_time`) VALUES (? , NOW())", online)
			}
		}
	}()
}

type tokenAdapter struct {
}

func (this tokenAdapter) Set(key string, value string) error {
	return ssdb.SSDB().Set(fmt.Sprintf("ALLHERO_%s", key), value)
}

func (this tokenAdapter) Get(key string) (string, error) {
	return ssdb.SSDB().Get(fmt.Sprintf("ALLHERO_%s", key))
}

func (this tokenAdapter) Delete(key string) error {
	return ssdb.SSDB().Del(fmt.Sprintf("ALLHERO_%s", key))
}
