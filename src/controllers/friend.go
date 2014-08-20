package controllers

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"models"
	"protodata"
)

func (this *Connect) FriendList() error {

	request := new(protodata.FriendListRequest)
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}

	uidList := models.User.UidList(request.GetSnsIds(), int(request.GetPlatId()))

	response := new(protodata.FriendListResponse)
	if len(uidList) == 0 {
		return this.Send(StatusOK, response)
	}

	FriendAction := models.NewFriendAction(this.Uid)

	var list1, list2 []*protodata.FriendData
	//	friendList := models.Role.FriendList(uidList)
	// --------- 临时写法 --------- //
	var friendList []*models.RoleData
	models.DB().Select(&friendList, "SELECT * FROM `role` ORDER BY `role_kill_num` DESC LIMIT 50")

	// --------- 临时写法 --------- //

	for index, f := range friendList {

		fdata := new(protodata.FriendData)
		fdata.Uid = proto.Int64(f.Uid)
		fdata.Num = proto.Int32(int32(index + 1))
		fdata.Point = proto.Int32(int32(f.KillNum))
		fdata.LeaderId = proto.Int32(int32(f.GeneralBaseId))
		fdata.Name = proto.String(f.Name)
		if _, ok := FriendAction.Map[f.Uid]; ok {
			fdata.IsGive = proto.Bool(true)
		} else {
			fdata.IsGive = proto.Bool(false)
		}
		list1 = append(list1, fdata)
	}

	// --------- 临时写法 --------- //
	friendList = make([]*models.RoleData, 0)
	models.DB().Select(&friendList, "SELECT * FROM `role` ORDER BY `role_unlimited_max_num` DESC LIMIT 50")

	// --------- 临时写法 --------- //
	//for i := 0; i < len(friendList); i++ {
	//	for j := len(friendList) - 1; j > i; j-- {
	//		if friendList[j].UnlimitedMaxNum > friendList[j-1].UnlimitedMaxNum {
	//			friendList[j-1], friendList[j] = friendList[j], friendList[j-1]
	//		}
	//	}
	//}

	for index, f := range friendList {
		fdata := new(protodata.FriendData)
		fdata.Uid = proto.Int64(f.Uid)
		fdata.Num = proto.Int32(int32(index + 1))
		fdata.Point = proto.Int32(int32(f.UnlimitedMaxNum))
		fdata.LeaderId = proto.Int32(int32(f.GeneralBaseId))
		fdata.Name = proto.String(f.Name)
		if _, ok := FriendAction.Map[f.Uid]; ok {
			fdata.IsGive = proto.Bool(true)
		} else {
			fdata.IsGive = proto.Bool(false)
		}
		list2 = append(list2, fdata)
	}

	response.FriendList1 = list1
	response.FriendList2 = list2
	return this.Send(StatusOK, response)
}

func (this *Connect) GiveAction() error {

	request := new(protodata.GiveStaminaRequest)
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}

	fid := request.GetUid()

	FriendAction := models.NewFriendAction(this.Uid)

	if len(FriendAction.Map) >= 20 {
		return this.Send(lineNum(), fmt.Errorf("今天已送了20个人不能再送"))
	}

	if _, ok := FriendAction.Map[fid]; ok {
		return this.Send(lineNum(), fmt.Errorf("今天已送过这个朋友不能再送"))
	}

	if err := FriendAction.Insert(fid); err != nil {
		return this.Send(lineNum(), err)
	}

	mail := new(models.MailData)
	mail.Uid = fid
	mail.ActionValue = 1
	mail.Title = "好友送体力"
	mail.Content = "好友:" + this.Role.Name + " 送你体力1点"

	if err := models.SendMail(mail); err != nil {
		return this.Send(lineNum(), err)
	}

	return this.Send(StatusOK, &protodata.GiveStaminaResponse{Uid: request.Uid})
}
