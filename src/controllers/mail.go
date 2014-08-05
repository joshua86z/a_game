package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"libs/lua"
	"models"
	"protodata"
)

func (this *Connect) MailList() error {

	var result []*protodata.MailData
	for _, mail := range models.NewMailModel(this.Role.Uid).List() {
		result = append(result, mailProto(mail))
	}

	Lua, _ := lua.NewLua("conf/notice.lua")
	content := Lua.GetString("content")
	datetime := Lua.GetString("datetime")
	Lua.Close()

	return this.Send(StatusOK, &protodata.MailResponse{
		Mails:        result,
		Cnnouncement: proto.String(content),
		Time:         proto.String(datetime)})
}

func (this *Connect) MailRewardRequest() error {

	request := &protodata.MailRewardRequest{}
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}

	mail := models.NewMailModel(this.Role.Uid).Mail(int(request.GetMailId()))
	if mail == nil {
		return this.Send(lineNum(), fmt.Errorf("没有这条邮件"))
	}

	models.DeleteMail(mail.Id)

	var rewardPoto protodata.RewardData

	if mail.ActionValue > 0 {
		if err := this.Role.SetActionValue(this.Role.ActionValue() + mail.ActionValue); err != nil {
			return this.Send(lineNum(), err)
		}
		rewardPoto.Stamina = proto.Int32(int32(mail.ActionValue))
	} else {
		if mail.Coin > 0 {
			this.Role.AddCoin(mail.Coin, models.FINANCE_MAIL_GET, fmt.Sprintf("mailId : %d", mail.Id))
			rewardPoto.RewardCoin = proto.Int32(int32(mail.Coin))
		}
		rewardPoto.RewardCoin = proto.Int32(int32(mail.Coin))
		if mail.Diamond > 0 {
			this.Role.AddDiamond(mail.Diamond, models.FINANCE_MAIL_GET, fmt.Sprintf("mailId : %d", mail.Id))
			rewardPoto.RewardDiamond = proto.Int32(int32(mail.Diamond))
		}
	}

	response := &protodata.MailRewardResponse{
		Role:   roleProto(this.Role),
		Reward: &rewardPoto}
	return this.Send(StatusOK, response)
}

func mailProto(mail *models.MailData) *protodata.MailData {

	reward := &protodata.RewardData{
		RewardCoin:    proto.Int32(int32(mail.Coin)),
		RewardDiamond: proto.Int32(int32(mail.Coin)),
		General:       nil,
		Stamina:       proto.Int32(int32(mail.ActionValue)),
	}
	return &protodata.MailData{
		MailId:      proto.Int32(int32(mail.Id)),
		MailTitle:   proto.String(mail.Title),
		MailContent: proto.String(mail.Content),
		Reward:      reward,
		IsReceive:   proto.Bool(mail.IsReceive)}

}
