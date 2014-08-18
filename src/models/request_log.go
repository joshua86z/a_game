package models

import (
	"libs/log"
	"time"
)

func init() {
	DB().AddTableWithName(requestLog{}, "request_logs").SetKeys(true, "Id")
}

type requestLog struct {
	Id      int   `db:"log_id"`
	Uid     int64 `db:"uid"`
	Index   int32 `db:"log_index"`
	AddTime int64 `db:"log_time"`
}

var (
	logChan chan *requestLog
)

func init() {
	logChan = make(chan *requestLog, 1000)
	go checkLogChan()
}

func checkLogChan() {

	defer func() {
		if err := recover(); err != nil {
			log.Critical("logChan panic : %v", err)
			checkLogChan()
		}
	}()

	for log := range logChan {
		DB().Insert(log)
	}
}

func InsertRequestLog(uid int64, index int32) {
	logChan <- &requestLog{
		Uid:     uid,
		Index:   index,
		AddTime: time.Now().Unix(),
	}
}
