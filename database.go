package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

func messageFromStoredReminder(stored StoredReminder) tb.Message {
	return tb.Message{ID: stored.MessageID, Chat: &tb.Chat{ID: stored.ChatID}}
}

func initDB(ctx context.Context) {
	dbCtx = ctx

	dbString := fmt.Sprintf("mongodb://%s:%s", mongoHostname, mongoPort)
	client, err := mongo.Connect(dbCtx, options.Client().ApplyURI(dbString))
	dbClient = client

	dbCursor = dbClient.Database("telegram")
	dbCol = dbCursor.Collection("reminders")

	if err != nil {
		log.Fatal("Failed to connect to database!")
	}
}

func getUserReminders(user *tb.User) ([]StoredReminder, error) {
	var reminders []StoredReminder

	cur, err := dbCol.Find(dbCtx, bson.M{"user.id": user.ID})
	defer cur.Close(dbCtx)

	if err != nil {
		return reminders, err
	}

	cur.All(dbCtx, &reminders)
	return reminders, nil
}

func loadStoredReminders() {
	var reminders []StoredReminder

	cur, err := dbCol.Find(dbCtx, bson.D{})
	defer cur.Close(dbCtx)

	if err != nil {
		log.Fatal("Error loading Saved Reminders")
	}
	cur.All(dbCtx, &reminders)

	for i := range reminders {
		timestamp := time.Unix(reminders[i].Timestamp, 0)
		duration := timestamp.Sub(time.Now())

		go forwardStoredMessageAfterDelay(reminders[i].ID, duration)
	}
}

func getStoredReminderFromID(id primitive.ObjectID) (StoredReminder, error) {
	reminder := StoredReminder{}

	res := dbCol.FindOne(dbCtx, bson.M{"_id": id})
	err := res.Decode(&reminder)

	if err != nil {
		log.Println("Unable to load DB message from ID")
		return reminder, err
	}

	return reminder, nil
}

func storeMessageIntoDB(m *tb.Message, recipient *tb.User, timestamp int64) primitive.ObjectID {
	// Returns ObjectID of stored document
	storedReminder := StoredReminder{ChatID: m.Chat.ID, MessageID: m.ID, User: recipient, Timestamp: timestamp}
	res, err := dbCol.InsertOne(dbCtx, storedReminder)

	if err != nil {
		log.Panic(err)
	}

	id := res.InsertedID.(primitive.ObjectID)
	return id
}

func removeMessageFromDB(id primitive.ObjectID) int64 {
	// Returns number of removed documents
	res, err := dbCol.DeleteMany(dbCtx, bson.M{"_id": id})

	if err != nil {
		log.Panic(err)
	}

	return res.DeletedCount
}

func removeAllUserMessages(*tb.User) error {
	_, err := dbCol.DeleteMany(dbCtx, bson.M{})

	if err != nil {
		return err
	}

	return nil
}
