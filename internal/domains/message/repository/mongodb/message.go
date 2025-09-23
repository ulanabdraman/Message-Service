package mongodb

import (
	"context"
	"fmt"
	"time"

	"MessageService/internal/domains/message/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageRepository interface {
	Insert(ctx context.Context, msg *model.Message) error
	GetByID(ctx context.Context, vehicleID int64) ([]*model.Message, error)
	GetByTimeRange(ctx context.Context, vehicleID int64, from, to time.Time) ([]*model.Message, error)
}

type messageRepository struct {
	db *mongo.Database
}

func NewMessageRepository(db *mongo.Database) MessageRepository {
	return &messageRepository{db: db}
}

func collectionName(vehicleID int64) string {
	return fmt.Sprintf("%d", vehicleID)
}

func (r *messageRepository) Insert(ctx context.Context, msg *model.Message) error {
	coll := r.db.Collection(collectionName(msg.ID))

	_, err := coll.InsertOne(ctx, bson.M{
		"dt":     msg.DT,
		"st":     msg.ST,
		"pos":    msg.Pos,
		"params": msg.Params,
	})
	return err
}

func (r *messageRepository) GetByID(ctx context.Context, vehicleID int64) ([]*model.Message, error) {
	coll := r.db.Collection(collectionName(vehicleID))

	cursor, err := coll.Find(ctx, bson.M{}) // без фильтра — все документы
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*model.Message
	for cursor.Next(ctx) {
		var msg model.Message
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *messageRepository) GetByTimeRange(ctx context.Context, vehicleID int64, from, to time.Time) ([]*model.Message, error) {
	coll := r.db.Collection(collectionName(vehicleID))

	filter := bson.M{
		"dt": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*model.Message
	for cursor.Next(ctx) {
		var msg model.Message
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
