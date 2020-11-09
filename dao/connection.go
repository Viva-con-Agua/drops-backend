package dao

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//DB represents the database client
var DB = new(mongo.Database)

//Connect connects the DB-Client with mongo database.
func Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	uri := "mongodb://" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("database connection failed", err)
	}
	DB = client.Database(os.Getenv("DB_NAME"))
}
