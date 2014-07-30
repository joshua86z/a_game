package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"models"
	"protodata"
)

func (this *Connect) MailList() (error) {

	var result []*protodata.MailData
	for _, mail := range models.NewMailModel(this.Role.Uid).List() {
		result = append(result, mailProto(mail))
	}

	return this.Send(StatusOK, &protodata.MailResponse{Mails: result})
}

func (this *Connect) MailRewardRequest() (error) {

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
	} else {
		if mail.Coin > 0 {
			rewardPoto.RewardCoin = proto.Int32(int32(mail.Coin))
			this.Role.AddCoin(mail.Coin, models.FINANCE_MAIL_GET, fmt.Sprintf("mailId : %d", mail.Id))
		}
		if mail.Diamond > 0 {
			rewardPoto.RewardDiamond = proto.Int32(int32(mail.Diamond))
			this.Role.AddDiamond(mail.Diamond, models.FINANCE_MAIL_GET, fmt.Sprintf("mailId : %d", mail.Id))
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
