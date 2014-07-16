package models

import (
	"fmt"
	"time"
)

// role_mail
type MailData struct {
	Id        int    `db:"mail_id"`
	Uid       int64  `db:"uid"`
	Title     string `db:"mail_title"`
	Content   string `db:"mail_content"`
	Reward    string `db:"mail_reward"`
	IsReceive bool   `db:"mail_is_receive"`
	UnixTime  int64  `db:"mail_time"`
}

func init() {
	DB().AddTableWithName(MailData{}, "role_mails").SetKeys(true, "Id")
}

type MailModel struct {
	Uid      int64
	MailList []*MailData
}

func NewMailModel(uid int64) *MailModel {

	var Mail MailModel

	var temp []*MailData
	_, err := DB().Select(&temp, "SELECT * FROM role_mails WHERE uid = ? ", uid)
	if err != nil {
		DBError(err)
	}

	Mail.Uid = uid
	Mail.MailList = temp

	return &Mail
}

func (this *MailModel) List() []*MailData {
	return this.MailList
}

func (this *MailModel) GetMail(mailId int) *MailData {
	for _, mail := range this.MailList {
		if mail.Id == mailId {
			return mail
		}
	}
	panic(fmt.Sprintf("没有这条数据 %d", mailId))
}

func InsertMail(mail *MailData) *MailData {

	mail.UnixTime = time.Now().Unix()
	if err := DB().Insert(mail); err != nil {
		DBError(err)
	}

	return mail
}

func DeleteMail(mailId int) error {
	_, err := DB().Exec("DELETE FROM role_mails WHERE mail_id = ? ", mailId)
	return err
}
