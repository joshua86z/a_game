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

	var find bool
	for index, f := range friendList {
		if f.Uid == this.Uid {
			find = true
		}
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

	if !find {
		fdata := new(protodata.FriendData)
		if this.Role.KillNum == friendList[len(friendList)-1].KillNum {
			*fdata.Num = *list1[len(list1)-1].Num + 1
		} else {
			sql := "SELECT COUNT(*) FROM `role` WHERE `role_kill_num` > ?"
			num, err := models.DB().SelectInt(sql, this.Role.KillNum)
			if err != nil {
				return this.Send(lineNum(), err)
			}
			fdata.Num = proto.Int32(int32(num + 1))
		}
		fdata.Uid = proto.Int64(this.Role.Uid)
		fdata.Point = proto.Int32(int32(this.Role.KillNum))
		fdata.LeaderId = proto.Int32(int32(this.Role.GeneralBaseId))
		fdata.Name = proto.String(this.Role.Name)
		fdata.IsGive = proto.Bool(true)
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

	find = false
	for index, f := range friendList {
		if f.Uid == this.Uid {
			find = true
		}
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

	if !find {
		fdata := new(protodata.FriendData)
		if this.Role.UnlimitedMaxNum == friendList[len(friendList)-1].UnlimitedMaxNum {
			*fdata.Num = *list1[len(list1)-1].Num + 1
		} else {
			sql := "SELECT COUNT(*) FROM `role` WHERE `role_unlimited_max_num` > ?"
			num, err := models.DB().SelectInt(sql, this.Role.UnlimitedMaxNum)
			if err != nil {
				return this.Send(lineNum(), err)
			}
			fdata.Num = proto.Int32(int32(num + 1))
		}
		fdata.Uid = proto.Int64(this.Role.Uid)
		fdata.Point = proto.Int32(int32(this.Role.UnlimitedMaxNum))
		fdata.LeaderId = proto.Int32(int32(this.Role.GeneralBaseId))
		fdata.Name = proto.String(this.Role.Name)
		fdata.IsGive = proto.Bool(true)
		list2 = append(list2, fdata)
	}

	response.FriendList1 = list1
	response.FriendList2 = list2
	response.GiveNum = proto.Int32(int32(len(FriendAction.Map)))
	response.GiveMax = proto.Int32(5)
	return this.Send(StatusOK, response)
}

func (this *Connect) GiveAction() error {

	request := new(protodata.GiveStaminaRequest)
	if err := Unmarshal(this.Request.GetSerializedString(), request); err != nil {
		return this.Send(lineNum(), err)
	}

	fid := request.GetUid()

	FriendAction := models.NewFriendAction(this.Uid)

	if len(FriendAction.Map) >= 5 {
		return this.Send(lineNum(), fmt.Errorf("今天已送了5个人不能再送"))
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
