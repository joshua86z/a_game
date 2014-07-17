package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"models"
	"protodata"
)

type Mail struct {
}

func (this *Mail) List(RoleModel *models.RoleModel, commandRequest *protodata.CommandRequest) (protodata.StatusCode, interface{}, error) {

	var result []*protodata.MailData
	for _, mail := range models.NewMailModel(RoleModel.Uid).List() {
		result = append(result, mailProto(mail))
	}

	return protodata.StatusCode_OK, &protodata.MailResponse{Mails: result}, nil
}

func (this *Mail) MailRewardRequest(RoleModel *models.RoleModel, commandRequest *protodata.CommandRequest) (protodata.StatusCode, interface{}, error) {

	request := &protodata.MailRewardRequest{}
	if err := Unmarshal(commandRequest.GetSerializedString(), request); err != nil {
		return lineNum(), nil, err
	}

	mail := models.NewMailModel(RoleModel.Uid).Mail(int(request.GetMailId()))
	if mail == nil {
		return lineNum(), nil, fmt.Errorf("没有这条邮件")
	}

	models.DeleteMail(mail.Id)

	var rewardPoto protodata.RewardData

	if mail.ActionValue > 0 {
		if err := RoleModel.SetActionValue(RoleModel.ActionValue() + mail.ActionValue); err != nil {
			return lineNum(), nil, err
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
	return protodata.StatusCode_OK, response, nil
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
