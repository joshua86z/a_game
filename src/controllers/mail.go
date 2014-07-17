package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"models"
	"protodata"
)

type Mail struct {
}

func (this *Mail) List(uid int64, commandRequest *protodata.CommandRequest) (string, error) {

	var result []*protodata.MailData
	for _, mail := range models.NewMailModel(uid).List() {
		result = append(result, mailProto(mail))
	}

	return ReturnStr(commandRequest, protodata.StatusCode_OK, &protodata.MailResponse{Mails: result}), nil
}

func (this *Mail) MailRewardRequest(uid int64, commandRequest *protodata.CommandRequest) (string, error) {

	request := &protodata.MailRewardRequest{}
	if err := Unmarshal(commandRequest.GetSerializedString(), request); err != nil {
		return ReturnStr(commandRequest, 26, ""), err
	}

	mail := models.NewMailModel(uid).Mail(int(request.GetMailId()))
	if mail == nil {
		return ReturnStr(commandRequest, 32, "没有这条邮件"), fmt.Errorf("没有这条邮件")
	}

	models.DeleteMail(mail.Id)

	var rewardPoto protodata.RewardData

	RoleModel := models.NewRoleModel(uid)
	if mail.ActionValue > 0 {
		if err := RoleModel.SetActionValue(RoleModel.ActionValue() + mail.ActionValue); err != nil {
			return ReturnStr(commandRequest, 46, "失败,数据库错误"), err
		}
	} else {
		if mail.Coin > 0 {
			rewardPoto.RewardCoin = proto.Int32(int32(mail.Coin))
			RoleModel.AddCoin(mail.Coin, models.MAIL_GET, fmt.Sprintf("mailId : %d", mail.Id))
		}
		if mail.Diamond > 0 {
			rewardPoto.RewardDiamond = proto.Int32(int32(mail.Diamond))
			RoleModel.AddDiamond(mail.Diamond, models.MAIL_GET, fmt.Sprintf("mailId : %d", mail.Id))
		}
	}

	response := &protodata.MailRewardResponse{
		Role:   roleProto(RoleModel),
		Reward: &rewardPoto}
	return ReturnStr(commandRequest, protodata.StatusCode_OK, response), nil
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
