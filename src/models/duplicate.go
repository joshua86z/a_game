package models

import (
	"fmt"
	"time"
)

// role_duplicates
type DuplicateData struct {
	Id       int    `db:"duplicate_id"`
	Uid      int64  `db:"uid"`
	Chapter int    `db:"duplicate_chapter"`
	Section     int `db:"duplicate_section"`
	UnixTime int64  `db:"duplicate_time"`
}

func init() {
	DB().AddTableWithName(DuplicateData{}, "role_duplicates").SetKeys(true, "Id")
}

type DuplicateModel struct {
	Uid         int64
	DuplicateList []*DuplicateData
}

func NewDuplicateModel(uid int64) *DuplicateModel {

	var Duplicate DuplicateModel

	var temp []*DuplicateData
	_, err := DB().Select(&temp, "SELECT * FROM role_duplicates WHERE uid = ?  ", uid)
	if err != nil {
		DBError(err)
	}

	Duplicate.Uid = uid
	Duplicate.DuplicateList = temp

	return &Duplicate
}

func (this *DuplicateModel) List() []*DuplicateData {
	return this.DuplicateList
}

func (this *DuplicateModel) GetDuplicate(duplicateId int) *DuplicateData {
	for _, duplicate := range this.DuplicateList {
		if duplicate.Id == duplicateId {
			return duplicate
		}
	}
	panic(fmt.Sprintf("没有这个关卡 %d", duplicateId))
}

func InsertDuplicate(uid int64, chapter int, section int) *DuplicateData {

	duplicate := &DuplicateData{}

	duplicate.Uid = uid
	duplicate.Chapter = chapter
	duplicate.Section = section
	duplicate.UnixTime = time.Now().Unix()

	if err := DB().Insert(duplicate); err != nil {
		DBError(err)
	}

	return duplicate
}

func DeleteDuplicate(duplicateId int) error {

	_, err := DB().Exec("DELETE FROM role_duplicates WHERE duplicate_id = ? ", duplicateId)
	return err
}
