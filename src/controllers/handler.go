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
	playerMap       map[int64]*Conn
	request_log_map map[int32]string
)

func init() {
	gameToken = token.NewToken(tokenAdapter{})
	playerMap = make(map[int64]*Conn)
	request_log_map = make(map[int32]string)
	CountOnline()
}

func Handler(ws *websocket.Conn) {

	online++
	// New Connect
	conn := &Conn{
		WsConn: ws,
		Chan:   make(chan string, 10),
	}

	conn.pushToClient()
	conn.PullFromClient()
	conn.Close()

	online--
}

// get one handler
func getHandler(index int32) func(int64, *protodata.CommandRequest) (string, error) {

	// return func
	if fun, ok := handlers[index]; ok {
		return fun
	}

	// return 404 func
	return func(uid int64, cq *protodata.CommandRequest) (string, error) {
		err := fmt.Errorf("No handler map to command:%d", cq.GetCmdId())
		return ReturnStr(cq, 999, ""), err
	}
}

func SendMessage(uid int64, message string) error {
	playLock.Lock()
	if conn, ok := playerMap[uid]; !ok {
		return fmt.Errorf("uid : %d not online", uid)
	} else {
		conn.Send(message)
	}
	playLock.Unlock()
	return nil
}

// Client Conn
type Conn struct {
	Uid    int64
	WsConn *websocket.Conn
	Chan   chan string
}

func (this *Conn) Send(s string) {
	this.Chan <- s
}

func (this *Conn) pushToClient() {
	go func() {
		for s := range this.Chan {
			var buf = make([]byte, 4)
			binary.LittleEndian.PutUint32(buf, uint32(len(s)))

			Buffer := bytes.NewBuffer(buf)
			Buffer.WriteString(s)

			if err := websocket.Message.Send(this.WsConn, Buffer.Bytes()); err != nil {
				log.Warn("Can't send msg. %v", err)
			} else {
				log.Info("Send Success")
			}
		}
	}()
}

func (this *Conn) InMap(uid int64) {
	this.Uid = uid
	playLock.Lock()
	if _, ok := playerMap[uid]; !ok {
		playerMap[uid] = this
	}
	playLock.Unlock()
}

func (this *Conn) Close() {
	playLock.Lock()
	delete(playerMap, this.Uid)
	playLock.Unlock()
	this.WsConn.Close()
}

// 从客户端读取信息
func (this *Conn) PullFromClient() {

	var uid int64
	var index int32
	for {

		// receive from ws conn
		var content string
		if err := websocket.Message.Receive(this.WsConn, &content); err != nil {
			if err.Error() == "EOF" {
				log.Info("Websocket receive EOF")
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

		var response string

		// parse proto message
		request, err := ParseContent(content)
		if err != nil {
			log.Error("Parse client request error. %v", err)
			response = ReturnStr(request, 9998, fmt.Sprintf("客户端错误:%v", err))
			continue
		}
		index = request.GetCmdId()

		// Panic recover
		defer func() {
			if err := recover(); err != nil {
				log.Critical("Panic occur. %v", err)
				response = ReturnStr(request, 9999, fmt.Sprintf("服务器错误:%v", err))
				this.Send(response)
			}
		}()

		if index != 10000 {
			// Check Login status
			if request.GetTokenStr() == "" {
				response = ReturnStr(request, protodata.StatusCode_INVALID_TOKEN, "")
			}
			uid, _ = gameToken.GetUid(request.GetTokenStr())
			if uid == 0 {
				response = ReturnStr(request, protodata.StatusCode_INVALID_TOKEN, "")
			} else {
				this.InMap(uid)
			}
		}

		// Checking true
		handler := getHandler(index)
		log.Info("Exec %v -> %s (uid:%d)", index, handlerNames[index], uid)

		// 执行命令
		if response, err = handler(uid, request); err != nil {
			log.Error("Exec command:%v error. %v", index, err)
		}

		execTime := time.Now().Sub(beginTime)
		if execTime.Seconds() > 0.1 {
			//慢日志
			log.Warn("Slow Exec -> %s, time is %v second", handlerNames[index], execTime.Seconds())
		} else {
			log.Info("time is %v second", execTime.Seconds())
		}

		// Send response to client
		this.Send(response)
		//this.Send <-

		if uid != 0 && index > 0 {
			//玩家操作记录
			if _, ok := request_log_map[index]; ok {
				//				requestLog.InsertLog(player.UniqueId, index)
			}
		}
	}
}

func (this *Conn) OtherRequest(request []byte) {

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
