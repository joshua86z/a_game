package controllers

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"libs/log"
	"models"
	"protodata"
	"strings"
	"time"
)

// Client Connect
type Connect struct {
	Role    *models.RoleModel
	Conn    *websocket.Connect
	Chan    chan string
	Request *protodata.CommectRequest
}

func (this *Connect) Send(code protodata.StatusCode, value interface{}) {
	this.Chan <- ReturnStr(this.Request, code, value)
	this.Request = nil
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
		// Panic recover
		defer func() {
			if err := recover(); err != nil {
				log.Critical("Panic occur. %v", err)
				this.Send(lineNum(), fmt.Sprintf("%v", err))
			}
		}()

		// receive from ws Connect
		var content string
		err := websocket.Message.Receive(this.Conn, &content)
		if err != nil {
			if err.Error() == "EOF" {
				log.Info("Conn receive EOF")
			} else {
				log.Error("Can't receive message. %v", err)
			}
			return
		}

		// **************** 其它接口 **************** //
		if strings.HasPrefix(content, "20140709_allhero_") {
			return this.OtherRequest([]byte(strings.Replace(content, "20140709_allhero_", "", len(content))))
		}
		// **************** 支付专用 **************** //

		beginTime := time.Now()
		log.Info(" Begin ")

		// parse proto message

		this.Request, err = ParseContent(content)
		if err != nil {
			log.Error("Parse client request error. %v", err)
			this.Send(lineNum(), fmt.Sprintf("%v", err))
			continue
		}

		if request.GetCmdId() != 10000 {
			// Check Login status
			if request.GetTokenStr() == "" {
				this.Send(protodata.StatusCode_INVALID_TOKEN, nil)
				continue
			}
			uid, _ := gameToken.GetUid(request.GetTokenStr())
			if uid == 0 {
				this.Send(protodata.StatusCode_INVALID_TOKEN, nil)
				continue
			} else {
				if this.Role == nil {
					this.Role = models.NewRoleModel(uid)
				} else if this.Role.Uid != uid {
					this.Send(protodata.StatusCode_INVALID_TOKEN, nil)
					continue
				}
			}
		}

		// 执行命令
		function := this.Function()
		function()
		log.Info("Exec -> %v (uid:%d)", function, RoleModel.Uid)

		execTime := time.Now().Sub(beginTime)
		if execTime.Seconds() > 0.1 {
			//慢日志
			log.Warn("Slow Exec -> %v, time is %v second", function, execTime.Seconds())
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
		this.Send(lineNum(), fmt.Sprintf("没有这方法 index : %d", index))
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
