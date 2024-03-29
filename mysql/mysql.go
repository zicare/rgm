package mysql

import (
	"database/sql"
	"fmt"
	"os"

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
		conn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
			cf.GetString("db.user"),
			cf.GetString("db.password"),
			cf.GetString("db.host"),
			cf.GetString("db.port"),
			cf.GetString("db.name"))
	)

	db, err = sql.Open("mysql", conn)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(cf.GetInt("db.max_open_conns"))

	return nil
}

//Db returns the db handler
func Db() *sql.DB {

	if db != nil {
		return db
	} else if err := Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return db
}
