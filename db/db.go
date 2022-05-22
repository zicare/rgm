package db

import (
	"database/sql"
	"fmt"

	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/msg"
)

var db *sql.DB

//Init tests the db connection and saves the db handler
func Init() error {

	var (
		err error
		c   = config.Config()
		//conn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			c.GetString("db.user"),
			c.GetString("db.password"),
			c.GetString("db.host"),
			c.GetString("db.port"),
			c.GetString("db.name"))
	)

	//db, err = sql.Open("postgres", conn)
	db, err = sql.Open("mysql", conn)
	if err != nil {
		//Server error: %s
		return msg.Get("25").SetArgs(err.Error()).M2E()
	}

	err = db.Ping()
	if err != nil {
		//Server error: %s
		return msg.Get("25").SetArgs(err.Error()).M2E()
	}

	db.SetMaxOpenConns(c.GetInt("db.max_open_conns"))

	return nil
}

//Db returns the db handler
func Db() *sql.DB {

	return db
}
