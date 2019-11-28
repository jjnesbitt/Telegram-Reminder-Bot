package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbClient *mongo.Client
	dbCtx    context.Context
	dbCursor *mongo.Database
	dbCol    *mongo.Collection
)

func initDB(ctx context.Context) {
	dbCtx = ctx
	client, err := mongo.Connect(dbCtx, options.Client().ApplyURI("mongodb://localhost:9000"))
	dbClient = client

	dbCursor = dbClient.Database("telegram")
	dbCol = dbCursor.Collection("reminders")

	if err != nil {
		log.Fatal("Failed to connect to database!")
	}
}
