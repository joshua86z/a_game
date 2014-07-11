package models

import (
	"fmt"
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

type ItemModel struct {
	Uid         int64
	ItemList []*ItemData
}

func GetItemModel(uid int64) *ItemModel {

	var Item ItemModel

	var temp []*ItemData
	_, err := DB().Select(&temp, "SELECT * FROM role_items WHERE uid = ?  ", uid)
	if err != nil {
		DBError(err)
	}

	Item.Uid = uid
	Item.ItemList = temp

	return &Item
}

func (this *ItemModel) List() []*ItemData {
	return this.ItemList
}

func (this *ItemModel) GetItem(itemId int) *ItemData {
	for _, item := range this.ItemList {
		if item.Id == itemId {
			return item
		}
	}
	panic(fmt.Sprintf("没有这个道具 %d", itemId))
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

func InsertItem(uid int64, configId int) *ItemData {

	c := ConfigItemMap()[configId]

	item := &ItemData{}

	item.Uid = uid
	item.ConfigId = configId
	item.Name = c.Name
	item.Level = 1
	item.UnixTime = time.Now().Unix()

	if err := DB().Insert(item); err != nil {
		DBError(err)
	}

	return item
}

func DeleteItem(itemId int) error {

	_, err := DB().Exec("DELETE FROM role_items WHERE item_id = ? ", itemId)
	return err
}
