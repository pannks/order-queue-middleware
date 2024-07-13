package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// fmt.Println("MONGODB_URI:", os.Getenv("MONGODB_URI"))
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)

	Client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// Ping the primary to verify the connection
	if err := Client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to MongoDB")
}
