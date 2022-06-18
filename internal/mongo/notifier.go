package mongo

import (
	"context"
	"github.com/Coflnet/auction-stats/internal/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func NotifiersForUser(userId int) ([]*model.Notifier, error) {
	filter := bson.D{{"user_id", userId}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := notifierCollection.Find(ctx, filter)

	if err != nil {
		log.Error().Err(err).Msgf("error listing notifiers, for user %d", userId)
		return nil, err
	}

	var notifiers []*model.Notifier
	for cur.Next(ctx) {
		var notifier model.Notifier
		err := cur.Decode(&notifier)
		if err != nil {
			log.Error().Err(err).Msgf("error decoding notifier")
			return nil, err
		}

		notifiers = append(notifiers, &notifier)
	}

	return notifiers, nil
}

func InsertNotifier(notifier *model.Notifier) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	notifier.ID = primitive.NewObjectID()

	res, err := notifierCollection.InsertOne(ctx, notifier)

	if err != nil {
		log.Error().Err(err).Msgf("error inserting notifier")
		return err
	}

	log.Info().Msgf("inserted notifier %v", res.InsertedID)
	return nil
}

func ReplaceNotifier(notifier *model.Notifier) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	res, err := notifierCollection.ReplaceOne(ctx, bson.D{{"_id", notifier.ID}}, notifier)

	if err != nil {
		log.Error().Err(err).Msgf("error replacing notifier")
		return err
	}

	log.Info().Msgf("replaced %d notifier", res.ModifiedCount)
	return nil
}

func DeleteNotifier(notifier *model.Notifier) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := notifierCollection.DeleteOne(ctx, bson.D{{"_id", notifier.ID}})
	if err != nil {
		log.Error().Err(err).Msgf("error deleting notifier")
		return err
	}

	log.Info().Msgf("deleted notifier %v", res.DeletedCount)
	return nil
}

func NotifierToEvaluate() (<-chan *model.Notifier, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"next_evaluation", bson.D{{"$lte", time.Now()}}}}

	cur, err := notifierCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	result := make(chan *model.Notifier)

	go func(c *mongo.Cursor) {
		for c.Next(ctx) {

			var notifier model.Notifier
			err := cur.Decode(&notifier)

			if err != nil {
				log.Error().Err(err).Msgf("error decoding notifier")
				return
			}

			result <- &notifier
		}

		close(result)
	}(cur)

	return result, nil
}
