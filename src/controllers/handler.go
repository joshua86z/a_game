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

const (
	LOGIN int32 = 10103
)

// Client Connect
type Connect struct {
	Uid     int64
	Role    *models.RoleModel
	Conn    *websocket.Conn
	Chan    chan string
	Request *protodata.CommandRequest
}

func (this *Connect) Send(code protodata.StatusCode, value interface{}) error {
	if _, ok := value.(error); ok {
		log.Error("Error: %v", value)
	}
	this.Chan <- ReturnStr(this.Request, code, value)
	this.Request = nil
	return nil
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

// 同一个账号只有一个连接
func (this *Connect) InMap(uid int64) {
	this.Uid = uid
	playLock.Lock()
	if val, ok := playerMap[this.Uid]; !ok {
		playerMap[this.Uid] = this
	} else {
		if this.Conn != val.Conn {
			val.Conn.Close()
			playerMap[this.Uid] = this
		}
	}
	playLock.Unlock()
}

func (this *Connect) Close() {
	playLock.Lock()
	if val, ok := playerMap[this.Uid]; ok {
		if this.Conn == val.Conn {
			delete(playerMap, this.Uid)
		}
	}
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
				log.Warn("Can't receive message. %v", err)
			}
			return
		}

		// **************** 其它接口 **************** //
		if strings.HasPrefix(content, "20140709_allhero_") {
			this.OtherRequest([]byte(strings.Replace(content, "20140709_allhero_", "", len(content))))
			return
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

		if this.Request.GetCmdId() != LOGIN {
			// Check Login status
			if this.Request.GetTokenStr() == "" {
				this.Send(protodata.StatusCode_INVALID_TOKEN, nil)
				continue
			}
			uid, _ := gameToken.GetUid(this.Request.GetTokenStr())
			if uid == 0 {
				this.Send(protodata.StatusCode_INVALID_TOKEN, nil)
				continue
			} else {
				if this.Role == nil {
					this.Role = models.NewRoleModel(uid)
					if this.Role == nil {
						this.Role = &models.RoleModel{}
						this.Role.Uid = uid
						this.Role.Coin = 0
						this.Role.Diamond = 0
						if err := models.InsertRole(this.Role); err != nil {
							this.Role = nil
							this.Send(lineNum(), err)
							continue
						}
					}

				} else if this.Role.Uid != uid {
					this.Send(protodata.StatusCode_INVALID_TOKEN, nil)
					continue
				}
			}
			this.InMap(uid)
		}

		if this.Uid > 0 {
			log.Info("Exec -> %d (uid:%d)", this.Request.GetCmdId(), this.Uid)
		}

		this.Function(this.Request.GetCmdId())()

		execTime := time.Now().Sub(beginTime)
		if int(execTime.Seconds()) > 1 {
			// slow log
			log.Warn("Slow Exec , time is %v second", execTime.Seconds())
		} else {
			log.Info("time is %v second", execTime.Seconds())
		}

		if this.Role != nil {
			// 玩家操作记录
			if _, ok := request_log_map[this.Request.GetCmdId()]; ok {
				//				requestLog.InsertLog(player.UniqueId, index)
			}
		}
	}
}

func (this *Connect) Function(index int32) func() error {
	switch index {
	case LOGIN:
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
		return this.BuyGeneral
	case 10112:
		return this.MailList
	case 10113:
		return this.MailRewardRequest
	case 10114:
		return this.ItemLevelUp
	default:
		return func() error {
			return this.Send(lineNum(), fmt.Sprintf("没有这方法 index : %d", index))
		}
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
