package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB = new(sql.DB)
var err = sql.ErrConnDone

func ConnectDatabase() {
	mysqlInfo := os.Getenv("DB_USER") +
		":" +
		os.Getenv("DB_PASSWORD") +
		"@tcp(" +
		os.Getenv("DB_HOST") +
		":" +
		os.Getenv("DB_PORT") +
		")/" +
		os.Getenv("DB_NAME")
	DB, err = sql.Open("mysql", mysqlInfo)
	if err != nil {
		log.Fatal(err, " ### utils.ConnectDatabase Step_1")
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal(err, " ### utils.ConnectDatabase Step_2")
	}
	fmt.Println("Database successfully connected!")
	//	query, err := ioutil.ReadFile("drops-database.sql")
	//	if err != nil {
	//		log.Fatal("read sql file failed", err)
	//	}
	//	if _, err := DB.Exec(string(query)); err != nil {
	//		log.Fatal("database init failed", err)
	//	}

}
