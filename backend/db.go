package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initDB() {
	MCLI, MCTX, MCAN = makeConnection()
}

func makeConnection() (*mongo.Client, context.Context, context.CancelFunc) {
	var username, password, database, server, port string

	username = "root"
	password = "mepco"
	// database = "customer"
	// server = "mongo"
	server = os.Getenv("MONGO_SERVER")
	port = "27017"

	connectionURI := fmt.Sprintf(connectionStringTemplate, username, password, server, port, database)

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))

	if err != nil {
		log.Println("Failed to create Client ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Println("Failed to connect with MongoDB ", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println("Failed to ping ", err)
	}

	log.Println("Connected to MongoDB")

	return client, ctx, cancel
}
