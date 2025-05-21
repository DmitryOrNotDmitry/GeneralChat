package infra

import (
	"context"
	"generalChat/entity"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatDB struct {
	client *mongo.Client
}

func CreateChatDB() *ChatDB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI("mongodb://mongo:mongo@localhost:27017/chatdb?authSource=admin")
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	return &ChatDB{client}
}

func (db *ChatDB) SaveMessage(message entity.Message) {
	messages := db.client.Database("chatdb").Collection("messages")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := messages.InsertOne(ctx, message)
	if err != nil {
		log.Fatal(err)
	}
}

func (db *ChatDB) GetLastNMessages(n int64) []entity.Message {
	messages := db.client.Database("chatdb").Collection("messages")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cursor, err := messages.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"timestamp": -1}).SetLimit(n))
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var result []entity.Message
	if err := cursor.All(ctx, &result); err != nil {
		log.Fatal(err)
	}

	return result
}

func (db *ChatDB) Close() error {
	return db.client.Disconnect(context.Background())
}
