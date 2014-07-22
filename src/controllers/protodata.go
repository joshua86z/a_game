package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"protodata"
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
func ReturnStr(r *protodata.CommandRequest, code protodata.StatusCode, obj interface{}) []byte {
	return Marshal(Return(r, &code, obj))
}

// 构造CommandResponse
func Return(r *protodata.CommandRequest, code *protodata.StatusCode, obj interface{}) *protodata.CommandResponse {

	var serStr, sstr *string
	sstr = proto.String("")
	switch msg := obj.(type) {
	case string:
		serStr = proto.String(msg)
		sstr = serStr
	case nil:
		serStr = proto.String("")
	case proto.Message:
		serStr = proto.String(string(Marshal(msg)))
	default:
		serStr = proto.String(fmt.Sprint(obj))
		//code = pb.StatusCode_SERVER_INTERNAL_ERROR
	}

	return &protodata.CommandResponse{
		Status:           &protodata.StatusData{SCode: code, SStr: sstr},
		CmdId:            proto.Int32(r.GetCmdId()),
		TokenStr:         proto.String(r.GetTokenStr()),
		CmdIndex:         proto.Int32(r.GetCmdIndex()),
		SerializedString: serStr,
	}
}

// 构造主动推送的CommandResponse
//func ActiveReturn(cmdId int32, code pb.StatusCode, obj interface{}) *pb.CommandResponse {
//
//	return nil
//}

// Parse ws receive content to CommandRequest
func ParseContent(content []byte) (*protodata.CommandRequest, error) {

	request := &protodata.CommandRequest{}
	return request, proto.Unmarshal(content, request)
}
