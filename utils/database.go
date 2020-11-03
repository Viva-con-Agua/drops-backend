package utils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB = new(sql.DB)
var err = sql.ErrConnDone

func ConnectDatabase() {
	mysqlInfo := os.Getenv("MYSQL_USER") +
		":" +
		os.Getenv("MYSQL_PASSWORD") +
		"@tcp(" +
		os.Getenv("DB_HOST") +
		":" +
		os.Getenv("DB_PORT") +
		")/" +
		os.Getenv("MYSQL_DATABASE")
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

var Database = new(mongo.Database)

func ConnectMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	uri := "mongodb://" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("database connection failed", err)
	}
	Database = client.Database("drops")
}
