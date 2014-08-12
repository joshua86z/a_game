package models

import (
	"database/sql"
	"time"
)

// role_items
type ItemData struct {
	Id       int    `db:"item_id"`
	Uid      int64  `db:"uid"`
	ConfigId int    `db:"item_config_id"`
	Name     string `db:"item_name"`
	Level    int    `db:"item_level"`
	UnixTime int64  `db:"item_time"`
}

func init() {
	DB().AddTableWithName(ItemData{}, "role_items").SetKeys(true, "Id")
}

var Item ItemModel

type ItemModel struct {
}

//func NewItemModel(uid int64) *ItemModel {

//	var Item ItemModel

//	var temp []*ItemData
//	_, err := DB().Select(&temp, "SELECT * FROM role_items WHERE uid = ? ", uid)
//	if err != nil {
//		DBError(err)
//	}

//	Item.Uid = uid
//	Item.ItemList = temp

//	return &Item
//}

func (this *ItemModel) List(uid int64) []*ItemData {
	var result []*ItemData
	_, err := DB().Select(&result, "SELECT * FROM role_items WHERE uid = ? ", uid)
	if err != nil {
		DBError(err)
	}
	return result
}

func (this ItemModel) Item(uid int64, configId int) *ItemData {

	ItemData := new(ItemData)
	err := DB().SelectOne(ItemData, "SELECT * FROM role_items WHERE uid = ? AND item_config_id = ?", uid, configId)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		DBError(err)
	}
	return ItemData
}

func (this ItemModel) Insert(uid int64, config *ConfigItem) *ItemData {

	item := new(ItemData)
	item.Uid = uid
	item.ConfigId = config.ConfigId
	item.Name = config.Name
	item.Level = 1
	item.UnixTime = time.Now().Unix()

	if err := DB().Insert(item); err != nil {
		return nil
	}

	return item
}

func (this *ItemData) LevelUp() error {

	this.Level += 1
	this.UnixTime = time.Now().Unix()

	_, err := DB().Update(this)
	return err
}

func (this *ItemData) LevelUpCoin() int {
	return this.Level * 10
}
