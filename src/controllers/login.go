package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"libs/log"
	"libs/lua"
	"models"
	"net/http"
	"protodata"
	"strconv"
	"strings"
	"time"
)

func (this *Connect) Login() error {

	request := &protodata.LoginRequest{}
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}

	platId := int(request.GetPlatId())
	username := request.GetUsername()
	password := request.GetPassword()
	otherId := request.GetOtherId()
	otherData := request.GetOtherData()
	session := request.GetOtherSession()
	sign := request.GetOtherSign()

	var user *models.UserData

	if platId == 2 { // PP助手
		otherId = ppLogin(session)
		if otherId == "" {
			return this.Send(lineNum(), fmt.Errorf("PP助手验证错误"))
		}
	} else {
		//else if platId == 4 { // 91助手
		//	if err := login91(otherId, session); err != nil {
		//		return this.Send(lineNum(), err)
		//	}
		//} else {
		var b bool
		otherId, b = otherLogin(platId, otherId, session, sign, otherData)
		if !b {
			return this.Send(lineNum(), fmt.Errorf("第三方验证错误"))
		}
	}

	user = models.User.GetUserByOtherId(otherId, platId)
	if user == nil {

		if platId == 0 {
			m := md5.New()
			m.Write([]byte(password))
			password = hex.EncodeToString(m.Sum(nil))
		}

		user = new(models.UserData)
		user.UserName = username
		user.Password = password
		user.OtherId = otherId
		user.Ip = request.GetIp()
		user.Imei = request.GetImei()
		user.PlatId = platId
		if err := user.Insert(); err != nil {
			return this.Send(lineNum(), err)
		}
	} else if platId == 0 {

		m := md5.New()
		m.Write([]byte(password))
		password = hex.EncodeToString(m.Sum(nil))
		if user.Password != password {
			return this.Send(lineNum(), fmt.Errorf("密码错误"))
		}
	}

	token, err := gameToken.AddToken(user.Uid)
	if err != nil {
		return this.Send(lineNum(), err)
	}

	log.Info("Exec -> login (uid:%d)", user.Uid)

	this.Uid = user.Uid
	if Role, err := models.Role.Role(user.Uid); err == nil {
		this.Role = Role
	}
	playerMap.Set(user.Uid, this)

	response := &protodata.LoginResponse{TokenStr: proto.String(token)}
	return this.Send(StatusOK, response)
}

func otherLogin(platId int, otherId, session, sign, otherData string) (string, bool) {

	if platId == 0 {
		return otherId, true
	}

	Lua, err := lua.NewLua(fmt.Sprintf("lua/plat_%d.lua", platId))
	if err != nil {
		log.Error("LUA ERROR : login.go line - 60")
		return "0", false
	}

	Lua.L.GetGlobal("isLogin")
	Lua.L.DoString(fmt.Sprintf("uid, isLogin = isLogin('%s', '%s', '%s', '%s')", otherId, session, sign, otherData))

	uid := Lua.GetString("uid")
	isLogin := Lua.GetBool("isLogin")

	Lua.Close()
	return uid, isLogin
}

func ppLogin(ppToken string) string {

	if len(ppToken) != 32 {
		return ""
	}

	pp_id := 4335
	app_key := "8dbbcdf221234073ccd75b1a277f7255"
	url := "http://passport_i.25pp.com:8080/account?tunnel-command=2852126760"

	m := md5.New()
	m.Write([]byte("sid=" + ppToken + app_key))
	sign := hex.EncodeToString(m.Sum(nil))

	postData := `{"data":{"sid":"%s"},"encrypt":"MD5","game":{"gameId":%d},"id":%d,"service":"account.verifySession","sign":"%s"}`
	postData = fmt.Sprintf(postData, ppToken, pp_id, time.Now().Unix(), sign)

	client := new(http.Client)
	resp, err := client.Post(url, "application/x-www-form-urlencoded", strings.NewReader(postData))
	if err != nil {
		log.Error("PPLOGIN ERROR , line: %d", lineNum())
		return ""
	}

	body, _ := ioutil.ReadAll(resp.Body)

	resp.Body.Close()

	var jsonData map[string]interface{}
	if json.Unmarshal(body, &jsonData) != nil {
		log.Error("PPLOGIN ERROR , line: %d", lineNum())
		return ""
	}

	code := jsonData["state"].(map[string]interface{})
	if int(code["code"].(float64)) != 1 {
		log.Error("PPLOGIN ERROR , line: %d", lineNum())
		return ""
	}

	data := jsonData["data"].(map[string]interface{})
	return data["accountId"].(string)
}

func login91(uin string, session string) error {

	appId := 100
	appKey := ""

	sign := strconv.Itoa(appId) + "4" + uin + session + appKey
	m := md5.New()
	m.Write([]byte(sign))
	sign = hex.EncodeToString(m.Sum(nil))

	url := "http://service.sj.91.com/usercenter/ap.aspx?AppId=%d&Act=4&Uin=%s&SessionId=%s&Sign=" + sign
	url = fmt.Sprintf(url, appId, uin, session)

	client := new(http.Client)
	resp, err := client.Get(url)
	if err != nil {
		log.Error("91LOGIN ERROR , line: %d , %v", lineNum(), err)
		return err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var jsonData map[string]string
	if err = json.Unmarshal(body, &jsonData); err != nil {
		log.Error("91LOGIN ERROR , line: %d , %v", lineNum(), err)
		return err
	}

	if jsonData["ErrorCode"] != "1" {
		err = fmt.Errorf("ErrorCode != 1")
		log.Error("91LOGIN ERROR , line: %d , %v", lineNum(), err)
		return err
	}

	return nil
}
