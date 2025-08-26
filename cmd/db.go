package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbClient struct {
	client         *mongo.Client
	dbName         string
	collectionName string
}

type DbUser struct {
	TelegramUserId int `bson:"telegram_user_id"`
	Preferences    []string
}

func LoadClient(ctx context.Context) (DbClient, error) {
	uri := os.Getenv("MONGO_URI")
	dn := os.Getenv("DATABASE_NAME")
	cn := os.Getenv("COLLECTION_NAME")
	dbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Printf("Error in connecting to MongoDB: %v", err)
		return DbClient{}, err
	}

	return DbClient{client: dbClient, dbName: dn, collectionName: cn}, err
}

func (db *DbClient) DisconnectClient(ctx context.Context) (*mongo.Client, error) {
	err := db.client.Disconnect(ctx)
	return db.client, err
}

func (db *DbClient) SaveUser(ctx context.Context, telegramUserId int) (DbUser, error) {
	coll := db.client.Database(db.dbName).Collection(db.collectionName)
	filter := bson.D{{Key: "telegram_user_id", Value: telegramUserId}}
	var result DbUser
	err := coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		// if user id not found, add it to db
		if err == mongo.ErrNoDocuments {
			log.Printf("New user: %v", telegramUserId)
			_, err := coll.InsertOne(ctx, DbUser{TelegramUserId: telegramUserId, Preferences: []string{}})
			if err != nil {
				log.Printf("Error in saving new user: %v", err)
				return DbUser{}, err
			}
			return DbUser{TelegramUserId: telegramUserId, Preferences: []string{}}, err
		} else {
			log.Printf("Error in querying the database: %v", err)
			return DbUser{}, err
		}
	}

	return result, err
}
