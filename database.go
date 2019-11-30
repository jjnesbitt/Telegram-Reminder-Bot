package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	dbClient *mongo.Client
	dbCtx    context.Context
	dbCursor *mongo.Database
	dbCol    *mongo.Collection
	dbCancel context.CancelFunc
)

func initDB(ctx context.Context) {
	dbCtx = ctx

	dbString := fmt.Sprintf("mongodb://%s:%s", mongoURL, mongoPort)
	client, err := mongo.Connect(dbCtx, options.Client().ApplyURI(dbString))
	dbClient = client

	dbCursor = dbClient.Database("telegram")
	dbCol = dbCursor.Collection("reminders")

	if err != nil {
		log.Fatal("Failed to connect to database!")
	}
}

func storeMessageIntoDB(m *tb.Message, recipient *tb.User, wait Wait) primitive.ObjectID {
	// Returns ObjectID of stored document
	storedReminder := StoredReminder{ChatID: m.Chat.ID, MessageID: m.ID, User: recipient, Time: wait.futureTimestamp}

	res, err := dbCol.InsertOne(dbCtx, storedReminder)

	if err != nil {
		log.Panic(err)
	}

	id := res.InsertedID.(primitive.ObjectID)
	return id
}

func removeMessageFromDB(id primitive.ObjectID) int64 {
	res, err := dbCol.DeleteMany(dbCtx, bson.M{"_id": id})

	if err != nil {
		log.Panic(err)
	}

	return res.DeletedCount
}
