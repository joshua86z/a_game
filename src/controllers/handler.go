package controllers

import (
	"code.google.com/p/go.net/websocket"
	"database/sql"
	"fmt"
	"libs/log"
	"models"
	"protodata"
	"time"
)

const (
	LOGIN int32 = 10103
)

// Client Connect
type Connect struct {
	Uid     int64
	Role    *models.RoleData
	Conn    *websocket.Conn
	Chan    chan []byte
	Request *protodata.CommandRequest
}

func (this *Connect) Send(status int, value interface{}) error {
	if _, ok := value.(error); ok {
		log.Error("linuNum -> %d : %v", status, value)
	}
	this.Chan <- ReturnStr(this.Request, status, value)
	//this.Request = nil
	return nil
}

func (this *Connect) pushToClient() {
	go func() {
		for s := range this.Chan {
			if _, err := this.Conn.Write(s); err != nil {
				log.Warn("Can't send msg. %v", err)
			} else {
				log.Info("Send Success")
			}
		}
	}()
}

func (this *Connect) Close() {
	this.Conn.Close()
	close(this.Chan)
}

// 从客户端读取信息
func (this *Connect) PullFromClient() {

	for {
		defer func() {
			if err := recover(); err != nil {
				log.Critical("Panic occur. %v", err)
				this.Send(lineNum(), fmt.Sprintf("%v", err))
				this.PullFromClient()
			}
		}()

		var content []byte
		err := websocket.Message.Receive(this.Conn, &content)
		if err != nil {
			log.Info("Websocket Error: %v", err)
			playerMap.Delete(this.Uid, this)
			return
		}

		beginTime := time.Now()
		log.Info(" Begin ")

		// parse proto message
		this.Request, err = ParseContent(content)
		if err != nil {
			log.Error("Parse client request error. %v", err)
			this.Send(lineNum(), err)
			continue
		}

		if this.Request.GetCmdId() != LOGIN {
			if !this.verify(this.Request.GetTokenStr()) {
				continue
			}
		}

		this.Function(this.Request.GetCmdId())()

		execTime := time.Now().Sub(beginTime)
		if execTime.Seconds() > 0.1 {
			// slow log
			log.Warn("Slow Exec , time is %v second", execTime.Seconds())
		} else {
			log.Info("time is %v second", execTime.Seconds())
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
	case 10109:
		return this.BuyDiamondRequest
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
	case 10115:
		return this.BattleRequest
	case 10116:
		return this.BattleResult
	case 10117:
		return this.SetLeader
	case 10118:
		return this.Sign
	case 10119:
		return this.FriendList
	case 10120:
		return this.GiveAction
	default:
		return func() error {
			return this.Send(lineNum(), fmt.Sprintf("没有这方法 index : %d", index))
		}
	}
}

func (this *Connect) verify(token string) bool {

	if token == "" {
		this.Send(2, nil)
		return false
	}
	uid, _ := gameToken.GetUid(token)
	if uid == 0 {
		this.Send(2, nil)
		return false
	}
	if this.Role == nil {
		this.Uid = uid
		Role, err := models.Role.Role(uid)
		if err == sql.ErrNoRows {
			if Role, err = models.Role.NewRole(uid); err != nil {
				// 插入失败
				return false
			}
		} else if err != nil {
			this.Send(lineNum(), err)
			return false
		}
		this.Role = Role
		playerMap.Set(uid, this)

	} else if this.Role.Uid != uid {
		this.Send(2, nil)
		return false
	} else {
		this.Role.UpdateDate()
	}

	log.Info("Exec -> %d (uid:%d)", this.Request.GetCmdId(), this.Uid)
	models.InsertRequestLog(this.Uid, this.Request.GetCmdId())
	return true
}
