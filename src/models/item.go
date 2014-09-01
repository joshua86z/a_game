package models

import (
	"database/sql"
	"time"
)

// role_items
type ItemData struct {
	Uid      int64  `db:"uid"`
	BaseId   int    `db:"item_base_id"`
	Name     string `db:"item_name"`
	Level    int    `db:"item_level"`
	UnixTime int64  `db:"item_time"`
}

func init() {
	DB().AddTableWithName(ItemData{}, "role_items").SetKeys(false, "Uid", "BaseId")
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

func (this ItemModel) Item(uid int64, baseId int) *ItemData {

	ItemData := new(ItemData)
	err := DB().SelectOne(ItemData, "SELECT * FROM role_items WHERE uid = ? AND item_base_id = ?", uid, baseId)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		DBError(err)
	}
	return ItemData
}

func (this ItemModel) Insert(uid int64, base *Base_Item) *ItemData {

	item := new(ItemData)
	item.Uid = uid
	item.BaseId = base.BaseId
	item.Name = base.Name
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
