package repo

import (
	"context"
	"errors"
	"generalChat/internal/model"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatDB struct {
	client   *mongo.Client
	database *mongo.Database
}

func CreateChatDB() *ChatDB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI("mongodb://mongo:mongo@localhost:27017/chatdb?authSource=admin")
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	db := client.Database("chatdb")
	return &ChatDB{client: client, database: db}
}

func (db *ChatDB) SaveMessage(message model.Message) error {
	if db == nil || db.database == nil {
		return errors.New("database is not initialized")
	}
	messages := db.database.Collection("messages")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := messages.InsertOne(ctx, message)
	if err != nil {
		return err
	}
	return nil
}

func (db *ChatDB) GetLastNMessages(n int64) ([]model.Message, error) {
	if db == nil || db.database == nil {
		return nil, errors.New("database is not initialized")
	}
	messages := db.database.Collection("messages")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cursor, err := messages.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"timestamp": -1}).SetLimit(n))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []model.Message
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (db *ChatDB) Close() error {
	if db == nil || db.client == nil {
		return nil
	}
	return db.client.Disconnect(context.Background())
}
