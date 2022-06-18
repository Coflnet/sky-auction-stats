package mongo

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

var (
	client *mongo.Client

	notifierCollection *mongo.Collection
)

func Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URL")))

	if err != nil {
		log.Error().Err(err).Msgf("error connecting to mongo")
		return err
	}

	db := client.Database("auction_stats")
	notifierCollection = db.Collection("notifiers")

	return nil
}

func Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return client.Disconnect(ctx)
}
