package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"libs/log"
	"libs/lua"
	"models"
	"protodata"
)

func (this *Connect) Login() error {

	request := &protodata.LoginRequest{}
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}

	platId := int(request.GetPlatId())
	ip := request.GetIp()
	username := request.GetUsername()
	password := request.GetPassword()
	otherId := request.GetOtherId()
	otherData := request.GetOtherData()
	session := request.GetOtherSession()
	sign := request.GetOtherSign()

	var user *models.UserData
	if otherId != "" && platId > 0 {
		otherId, b := otherLogin(platId, otherId, session, sign, otherData)
		if !b {
			return this.Send(lineNum(), fmt.Errorf("第三方验证错误"))
		}
		user = models.User.GetUserByOtherId(otherId, platId)
	} else {
		m := md5.New()
		m.Write([]byte(password))
		password = hex.EncodeToString(m.Sum(nil))
		user = models.User.GetUserByName(username)
		if user != nil && user.Password != password {
			return this.Send(lineNum(), fmt.Errorf("密码错误"))
		}
	}

	if user == nil {
		user = new(models.UserData)
		user.UserName = username
		user.Password = password
		user.OtherId = otherId
		user.Ip = ip
		user.PlatId = platId
		if err := user.Insert(); err != nil {
			return this.Send(lineNum(), err)
		}
	}

	token, err := gameToken.AddToken(user.Uid)
	if err != nil {
		return this.Send(lineNum(), err)
	}

	log.Info("Exec -> login (uid:%d)", user.Uid)
	this.InMap(user.Uid)
	this.Role = models.Role.Role(user.Uid)

	response := &protodata.LoginResponse{TokenStr: proto.String(token)}
	return this.Send(StatusOK, response)
}

func otherLogin(platId int, otherId, session, sign, otherData string) (string, bool) {

	Lua, err := lua.NewLua(fmt.Sprintf("lua/plat_%d.lua", platId))
	if err != nil {
		log.Error("LUA ERROR : login.go line - 60")
		return "0", false
	}

	//	Lua.L.GetGlobal("getUrl")
	//	Lua.L.DoString(fmt.Sprintf("url = getUrl('%s', '%s', '%s', '%s')", otherId, session, sign, otherData))
	//
	//	requestUrl := Lua.GetString("url")
	//
	//	Lua.L.GetGlobal("isPost")
	//	Lua.L.DoString("post = isPost()")
	//
	//	isPost := Lua.GetBool("post")
	//
	//	method := "GET"
	//	reader := &strings.Reader{}
	//	if isPost {
	//		Lua.L.GetGlobal("getPost")
	//		Lua.L.DoString(fmt.Sprintf("postData = getPost('%s', '%s', '%s', '%s')", otherId, session, sign, otherData))
	//		reader = strings.NewReader(Lua.GetString("postData"))
	//		method = "POST"
	//	}
	//
	//	request, err := http.NewRequest(method, requestUrl, reader)
	//	if err != nil {
	//		return false
	//	}
	//
	//	client := &http.Client{}
	//	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	//	request.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
	//	request.Header.Set("Accept-Language", "en-US,en;q=0.8,zh-CN;q=0.6,zh;q=0.4,ja;q=0.2")
	//	request.Header.Set("Cache-Control", "max-age=0")
	//	request.Header.Set("Connection", "keep-alive")
	//	if method == "POST" {
	//		request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	//	}
	//	result, err := client.Do(request)
	//	if err != nil {
	//		log.Error("HTTP REQUEST ERROR : login.go line - 98")
	//		return false
	//	}
	//
	//	body, _ := ioutil.ReadAll(result.Body)
	//	result.Body.Close()

	Lua.L.GetGlobal("isLogin")
	Lua.L.DoString(fmt.Sprintf("uid, isLogin = isLogin('%s', '%s', '%s', '%s')", otherId, session, sign, otherData))

	uid := Lua.GetString("uid")
	isLogin := Lua.GetBool("isLogin")

	Lua.Close()
	return uid, isLogin
}
