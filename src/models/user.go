package models

import (
	"database/sql"
	"fmt"
	"time"
)

func init() {
	var UserModel UserModel
	DB().AddTableWithName(UserModel, UserModel.tableName()).SetKeys(true, "Uid")
}

// user
type UserModel struct {
	Uid      int64  `db:"uid"`
	UserName string `db:"username"`
	Password string `db:"password"`
	Ip       string `db:"ip"`
	OtherId  string `db:"other_id"`
	PlatId   int    `db:"plat_id"`
	RegTime  int64  `db:"reg_time"`
}

func (this UserModel) tableName() string {
	return "user"
}

func (this UserModel) Insert() error {
	this.RegTime = time.Now().Unix()
	return DB().Insert(&this)
}

func GetUserByName(name string) *UserModel {

	UserModel := new(UserModel)

	str := "SELECT * FROM " + UserModel.tableName() + " WHERE username = ? LIMIT 1"
	if err := DB().SelectOne(UserModel, str, name); err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			panic(err)
		}
	}

	return UserModel
}

func NewUserModel(uid int64) *UserModel {

	UserModel := new(UserModel)

	err := DB().SelectOne(UserModel, "SELECT * FROM "+UserModel.tableName()+" WHERE uid = ? LIMIT 1", uid)
	if err != nil {
		panic(fmt.Sprintf("NewUserModel Error : %v", err))
	}

	return UserModel
}

func GetUserByOtherId(otherId string, platId int) *UserModel {

	UserModel := new(UserModel)

	err := DB().SelectOne(UserModel, "SELECT * FROM "+UserModel.tableName()+" WHERE other_id = ? AND plat_id = ? LIMIT 1", otherId, platId)
	if err == sql.ErrNoRows {
		return nil
	} else {
		panic(err)
	}

	return UserModel
}
