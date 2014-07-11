package db

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"libs/lua"
)

var db *gorp.DbMap

func init() {

	setting, err := lua.NewLua("conf/app.lua")
	if err != nil {
		panic(err)
	}

	db = Open(setting.GetString("dsn"))

	setting.Close()
}

func Open(dsn string) *gorp.DbMap {

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	return &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
}

func DB() *gorp.DbMap {
	return db
}
