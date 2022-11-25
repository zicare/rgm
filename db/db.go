package db

import (
	"database/sql"
	"fmt"

	"github.com/zicare/rgm/config"

	//required
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

//Init tests the db connection and saves the db handler
func Init() error {

	var (
		err  error
		cf   = config.Config()
		conn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			cf.GetString("db.user"),
			cf.GetString("db.password"),
			cf.GetString("db.host"),
			cf.GetString("db.port"),
			cf.GetString("db.name"))
	)

	db, err = sql.Open("mysql", conn)
	if err != nil {
		return new(OpenConnError)
	}

	err = db.Ping()
	if err != nil {
		return new(PingTestError)
	}

	db.SetMaxOpenConns(cf.GetInt("db.max_open_conns"))

	return nil
}

//Db returns the db handler
func Db() *sql.DB {

	return db
}
