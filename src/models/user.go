package models

import (
	"database/sql"
	"strings"
	"time"
)

func init() {
	User.tableName = "user"
	DB().AddTableWithName(UserData{}, User.tableName).SetKeys(true, "Uid")
}

var User UserModel

type UserModel struct {
	tableName string
}

// user
type UserData struct {
	Uid      int64  `db:"uid"`
	UserName string `db:"username"`
	Password string `db:"password"`
	OtherId  string `db:"other_id"`
	Ip       string `db:"ip"`
	Imei     string `db:"imei"`
	PlatId   int    `db:"plat_id"`
	RegTime  int64  `db:"reg_time"`
}

func (this *UserData) Insert() error {
	this.RegTime = time.Now().Unix()
	return DB().Insert(this)
}

func (this UserModel) GetUserByName(name string) *UserData {

	UserData := new(UserData)

	str := "SELECT * FROM " + User.tableName + " WHERE username = ? LIMIT 1"
	if err := DB().SelectOne(UserData, str, name); err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else if err != nil {
			DBError(err)
		}
	}

	return UserData
}

func (this UserModel) User(uid int64) *UserData {

	UserData := new(UserData)

	err := DB().SelectOne(UserData, "SELECT * FROM "+User.tableName+" WHERE uid = ? LIMIT 1", uid)
	if err != nil {
		DBError(err)
	}

	return UserData
}

func (this UserModel) GetUserByOtherId(otherId string, platId int) *UserData {

	UserData := new(UserData)

	err := DB().SelectOne(UserData, "SELECT * FROM "+User.tableName+" WHERE other_id = ? AND plat_id = ? LIMIT 1", otherId, platId)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		DBError(err)
	}

	return UserData
}

func (this UserModel) UidList(otherIds []string, platId int) []int64 {

	var temp []*UserData

	DB().Select(&temp, "SELECT * FROM "+User.tableName+" WHERE plat_id = ? AND other_id IN('"+strings.Join(otherIds, "','")+"')", platId)

	uids := make([]int64, len(temp))
	for index, user := range temp {
		uids[index] = user.Uid
	}

	return uids
}
