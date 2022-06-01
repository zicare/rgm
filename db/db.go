package db

import (
	"database/sql"
	"fmt"

	"github.com/zicare/rgm/config"
	"github.com/zicare/rgm/msg"

	//required
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var db *sql.DB

//Init tests the db connection and saves the db handler
func Init(flavor string) error {

	var (
		err  error
		conn string
		c    = config.Config()
	)

	switch flavor {
	case "mysql":
		conn = "%s:%s@tcp(%s:%s)/%s?parseTime=true"
	case "postgres":
		conn = "postgres://%s:%s@%s:%s/%s?sslmode=disable"
	default:
		return msg.Get("25").SetArgs(err.Error()).M2E()
	}

	conn = fmt.Sprintf(conn,
		c.GetString("db.user"),
		c.GetString("db.password"),
		c.GetString("db.host"),
		c.GetString("db.port"),
		c.GetString("db.name"))

	db, err = sql.Open(flavor, conn)
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
