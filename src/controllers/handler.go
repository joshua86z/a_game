package controllers

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	//	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"libs/log"
	"libs/ssdb"
	"libs/token"
	"strings"
	"time"
	"encoding/binary"
	//requestLog "game/models/logs"
	"encoding/json"
	"protodata"
)

// 运行变量
var (
	FullStaminaNum  int = 20 // 满体力数值
	gameToken       *token.Token
	request_log_map map[int32]string
)

func init() {

	gameToken = token.NewToken(tokenAdapter{})
	request_log_map = make(map[int32]string)
}

// Client Conn
type Conn struct {
	WsConn     *websocket.Conn
	onlineTime time.Time
	lastTime   time.Time
	//	Send   chan string
}

func (this *Conn) Send(s string) {

	var buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(len(s)))

	Buffer := bytes.NewBuffer(buf)
	Buffer.WriteString(s)

	if err := websocket.Message.Send(this.WsConn, Buffer.Bytes()); err != nil {
		log.Error("Can't send msg. %v", err)
	} else {
		log.Info("Send Success")
	}
}

// 发送信息回客户端
//func (this *Conn) PushToClient() {
//
//	// range Send chan
//	for s := range this.Send {
//
//		// response bytes 前4字节代表长度
//		//		log.Debug("Send LEN: %d", len(s))
//		buf := bytes.NewBuffer(common.IntToBytes(len(s)))
//		buf.WriteString(s)
//
//		if err := websocket.Message.Send(this.WsConn, buf.Bytes()); err != nil {
//			log.Error("Can't send msg. %v", err)
//		} else {
//			log.Info("Send Success")
//		}
//	}
//}

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

		this.lastTime = time.Now()
		//		if int(this.lastTime.Sub(this.onlineTime).Seconds()) >= task.OnlineTimeInterval {
		//			this.onlineTime = this.lastTime
		//			task.OnlineAddNum(this.lastTime)
		//		}

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

//
//// 验证玩家
//func (this *Conn) getPlayer(cq *protodata.CommandRequest) (*Player, error) {
//
//	token := cq.GetTokenStr()
//
//	log.Info("Token : %s", token)
//
//	uniqueId, err := Token().GetUid(token)
//	//	uniqueId, err := models.GetUnique(token)
//	if err != nil {
//		return nil, err
//	}
//
//	// token invalid or expires
//	if uniqueId == 0 {
//		log.Info("Token:%s invalid or Expired", token)
//		return nil, fmt.Errorf(fmt.Sprintf("Token:%s invalid or Expired", token))
//	}
//
//	return &Player{
//		Uid:   uniqueId,
//		Token:      token,
//		Conn:     this.WsConn,
//		ActiveTime: time.Now().Unix(),
//	}, nil
//}

// 重写 Handler
func Handler(ws *websocket.Conn) {

	onlineTime := time.Now().Add(-time.Second * 5 * 60)
	// New Conn
	conn := &Conn{
		WsConn:     ws,
		onlineTime: onlineTime,
		lastTime:   onlineTime,
		//		Send:   make(chan string, 3),
	}

	// Handle Websocket Message
	//	go conn.PushToClient()
	conn.PullFromClient()

	// kill resource
	conn.WsConn.Close()
	//	close(conn.Send)
}

// 10107 选择游戏服务器
//func selectServer(cq *protodata.CommandRequest) (string, error) {
//
//	log.Info("Exec 10107 -> selectServer")
//
//	request := &protodata.SelectServerRequest{}
//	err := proto.Unmarshal([]byte(cq.GetSerializedString()), request)
//	if err != nil {
//		return pd.ReturnStr(cq, protodata.StatusCode_DATA_ERROR, ""), err
//	}
//
//	uid64, err := token.NewToken(token.GateWayTokenAdapter()).GetUid(request.GetGwTokenStr())
//	uid := int(uid64)
//
//	// gwServer未登录
//	if uid == 0 {
//		return pd.ReturnStr(cq, protodata.StatusCode_INVALID_TOKEN,
//			&protodata.SelectServerResponse{
//				ECode: proto.Int32(2),
//				EStr:  proto.String(common.ESTR_token_expired),
//			}), nil
//	}
//
//	// 判断是否存在角色
//	userRole, err := models.GetRoleByUid(uid)
//	if err != nil {
//		return pd.ReturnStr(cq, protodata.StatusCode_DATABASE_ERROR, ""), err
//	}
//
//	// 如果不存在则生成角色UniqueId
//	var flag bool
//	var lockFlag bool
//	var uniqueId int64
//	if userRole == nil {
//
//		uniqueId, err = strconv.ParseInt(fmt.Sprintf("%d%d", 100+ServerId, uid), 10, 64)
//		if err != nil {
//			panic(err)
//		}
//	} else {
//		flag = true
//		uniqueId = userRole.Unique
//		if userRole.Status != 1 {
//			lockFlag = true
//		}
//	}
//
//	// 账户未被锁定
//	var serverToken string
//	if !lockFlag {
//
//		serverToken, err = Token().AddToken(uniqueId)
//		if err != nil {
//			panic(err)
//		}
//
//		// 设置最近登录服务器
//		ssdb.SSDB().Set(fmt.Sprintf("uid:%d:recent", uid), strconv.Itoa(int(ServerId)))
//	}
//
//	response := &protodata.SelectServerResponse{
//		ECode:          proto.Int32(1),
//		ServerTokenStr: proto.String(serverToken),
//		RoleExistFlag:  proto.Bool(flag),
//		RoleLockFlag:   proto.Bool(lockFlag),
//	}
//
//	return pd.ReturnStr(cq, protodata.StatusCode_OK, response), nil
//}

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

func (this *Conn) OtherRequest(request []byte) {

	data := make(map[string]string)

	json.Unmarshal(request, &data)

	if data["cmd"] == "8889" {
		this.Send("")
	}

	//	if data["cmd"] == "9000" {
	//		if err := Pay.ConfirmOrder(data["orderId"]); err != nil {
	//			log.Error("Pay Error %v", err)
	//		}
	//	} else if data["cmd"] == "9010" {
	//		configs.ConfigRefreshGeneralReset()
	//	}
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
