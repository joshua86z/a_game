package models

import (
	"libs/db"
	"github.com/coopernurse/gorp"
	"fmt"
)

func init() {
}

func DB() *gorp.DbMap {
	return db.DB
}

func DBError(err error) {
	panic(fmt.Sprintf("数据库错误: %v", err))
}
