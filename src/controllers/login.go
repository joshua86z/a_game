package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"models"
	"protodata"
//	"net/http"
//	"io/ioutil"
	"libs/lua"
//	"strings"
	"libs/log"
)

type Login struct {
}

func (this *Login) Login(uid int64, commandRequest *protodata.CommandRequest) (string, error) {

	request := &protodata.LoginRequest{}
	if err := Unmarshal(commandRequest.GetSerializedString(), request); err != nil {
		return ReturnStr(commandRequest, 19, ""), fmt.Errorf("%v", err)
	}

	username := request.GetUsername()
	password := request.GetPassword()
	otherId := request.GetOtherId()
	platId := int(request.GetPlatId())

	UserModel := models.GetUserByName(username)
	if UserModel.Uid == 0 {
		UserModel.UserName = username
		UserModel.Password = password
		UserModel.OtherId = otherId
		UserModel.PlatId = platId
		if err := UserModel.Insert(); err != nil {
			return ReturnStr(commandRequest, 34, ""), err
		}
	}

	uid = UserModel.Uid
	token, err := gameToken.AddToken(uid)
	if err != nil {
		return ReturnStr(commandRequest, 41, ""), fmt.Errorf("%v", err)
	}

	models.NewRoleModel(uid)
	models.NewSignModel(uid)

	response := &protodata.LoginResponse{TokenStr: proto.String(token)}
	return ReturnStr(commandRequest, 1, response), nil
}

func otherLogin(platId int, otherId string, session string, sign string, otherData string) (string, bool) {

	Lua, err := lua.NewLua(fmt.Sprintf("lua/%d.lua", platId))
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
