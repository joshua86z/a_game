package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"protodata"
)

const (
	StatusOK int = 1
)

// Marshal proto，保证msg都是正确的，否则panic
func Marshal(msg proto.Message) []byte {
	b, err := proto.Marshal(msg)
	if err != nil {
		panic(err.Error())
	}
	return b
}

// 封装反序列化proto (or 直接使用proto.Unmarshal?)
func Unmarshal(str string, msg proto.Message) error {
	return proto.Unmarshal([]byte(str), msg)
}

// 构造CommandResponse并Marshal成字符串
func ReturnStr(r *protodata.CommandRequest, status int, obj interface{}) []byte {
	return Marshal(Return(r, status, obj))
}

// 构造CommandResponse
func Return(r *protodata.CommandRequest, status int, obj interface{}) *protodata.CommandResponse {

	var serStr *string
	var code protodata.StatusData
	code.SCode = proto.Int32(int32(status))
	switch msg := obj.(type) {
	case string:
		serStr = proto.String(msg)
		code.SStr = serStr
	case nil:
		serStr = proto.String("")
	case proto.Message:
		serStr = proto.String(string(Marshal(msg)))
	default:
		serStr = proto.String(fmt.Sprint(obj))
		code.SStr = serStr
	}

	return &protodata.CommandResponse{
		Status:           &code,
		CmdId:            proto.Int32(r.GetCmdId()),
		TokenStr:         proto.String(r.GetTokenStr()),
		CmdIndex:         proto.Int32(r.GetCmdIndex()),
		SerializedString: serStr,
	}
}

// Parse ws receive content to CommandRequest
func ParseContent(content []byte) (*protodata.CommandRequest, error) {
	request := &protodata.CommandRequest{}
	return request, proto.Unmarshal(content, request)
}
