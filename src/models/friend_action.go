package models

import (
	"time"
)

func init() {
	DB().AddTableWithName(FriendActionModel{}, "friend_action").SetKeys(false, "Uid" , "Fid", "Date")
}

type FriendAction struct {
	Uid int64
	Map map[int64]*FriendActionModel
}

// friend_action
type FriendActionModel struct {
	Uid  int64  `db:"uid"`
	Fid  int64  `db:"friend_uid"`
	Date string `db:"date"`
}

func NewFriendAction(uid int64) *FriendAction {

	var temp []*FriendActionModel
	DB().Select(&temp, "SELECT * FROM friend_action WHERE uid = ? AND `date` = ?", uid, time.Now().Format("20060102"))

	FriendAction := new(FriendAction)
	FriendAction.Uid = uid
	FriendAction.Map = make(map[int64]*FriendActionModel)

	for _, val := range temp {
		FriendAction.Map[val.Fid] = val
	}

	return FriendAction
}

func (this FriendAction) Insert(fid int64) error {
	FriendActionModel := new(FriendActionModel)
	FriendActionModel.Uid = this.Uid
	FriendActionModel.Fid = fid
	FriendActionModel.Date = time.Now().Format("20060102")
	err := DB().Insert(FriendActionModel)
	if err != nil {
		this.Map[fid] = FriendActionModel
	}
	return err
}
