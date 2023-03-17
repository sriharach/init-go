package configs

import (
	"api-enjor/pkg/utils"
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	dsn, _ := utils.ConnectionURLBuilder("POSTGRES")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database")
	}

	fmt.Println("Database connected!")

	return db
}

func InintMongodbAtlas() *mongo.Client {
	// Set up client options
	clientOptions := options.Client().ApplyURI("mongodb+srv://" + os.Getenv("MONGODB_USERNAME") + ":<" + os.Getenv("MONGODB_PASSWORD") + ">@cluster0.re7aemm.mongodb.net/test")

	// Connect to MongoDB Atlas
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic("Failed to connect to MongoDB Atlas:" + err.Error())
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic("Failed to ping MongoDB Atlas:" + err.Error())
	}

	fmt.Println("Connected to MongoDB Atlas!")

	return client
}
