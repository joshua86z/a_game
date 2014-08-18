package models

import (
	"strconv"
	"time"
)

// role_duplicates
type DuplicateData struct {
	Id       int   `db:"duplicate_id"`
	Uid      int64 `db:"uid"`
	Chapter  int   `db:"duplicate_chapter"`
	Section  int   `db:"duplicate_section"`
	UnixTime int64 `db:"duplicate_time"`
}

func init() {
	DB().AddTableWithName(DuplicateData{}, "role_duplicates").SetKeys(true, "Id")
}

type DuplicateModel struct {
	Uid           int64
	DuplicateList []*DuplicateData
}

func NewDuplicateModel(uid int64) *DuplicateModel {

	Duplicate := new(DuplicateModel)

	var temp []*DuplicateData
	_, err := DB().Select(&temp, "SELECT * FROM role_duplicates WHERE uid = ? ", uid)
	if err != nil {
		DBError(err)
	}

	Duplicate.Uid = uid
	Duplicate.DuplicateList = temp

	return Duplicate
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
	return nil
}

func (this *DuplicateModel) Insert(chapter int, section int) *DuplicateData {

	duplicate := new(DuplicateData)

	duplicate.Uid = this.Uid
	duplicate.Chapter = chapter
	duplicate.Section = section
	duplicate.UnixTime = time.Now().Unix()

	if err := DB().Insert(duplicate); err != nil {
		return nil
	} else {
		this.DuplicateList = append(this.DuplicateList, duplicate)
	}

	return duplicate
}

func DuplicateTop(uids []int64) ([]int64, []int) {

	var temp []struct {
		Uid int64 `db:"uid"`
		Num int   `db:"num"`
	}

	var s string
	for _, uid := range uids {
		s += strconv.Itoa(int(uid)) + ","
	}
	sql := "SELECT * FROM (SELECT COUNT(*) AS num,uid FROM role_duplicates WHERE uid IN(" + s[0:len(s)-1] + ") GROUP BY uid) AS t ORDER BY num"
	DB().Select(&temp, sql)

	uidList := make([]int64, len(temp))
	numList := make([]int, len(temp))
	for _, val := range temp {
		uidList = append(uidList, val.Uid)
		numList = append(numList, val.Num)
	}

	return uidList, numList
}
