package mongodb

import (
	"context"
	"log"
	"time"

	"MessageService/internal/domains/message/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepository interface {
	Insert(ctx context.Context, msg *model.Message) error
	GetByID(ctx context.Context, id int64) (*model.Message, error)
	GetByTimeRange(ctx context.Context, id int64, from, to time.Time) ([]*model.Message, error)
}

type messageRepository struct {
	collection *mongo.Collection
}

func NewMessageRepository(db *mongo.Database, collectionName string) MessageRepository {
	return &messageRepository{
		collection: db.Collection(collectionName),
	}
}

// Insert сообщение
func (r *messageRepository) Insert(ctx context.Context, msg *model.Message) error {
	_, err := r.collection.InsertOne(ctx, msg)
	return err
}

// GetByID по ObjectID
func (r *messageRepository) GetByID(ctx context.Context, id int64) (*model.Message, error) {
	var msg model.Message
	err := r.collection.FindOne(ctx, bson.M{"uuid": id}).Decode(&msg)
	log.Println(msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *messageRepository) GetByTimeRange(ctx context.Context, id int64, from, to time.Time) ([]*model.Message, error) {
	filter := bson.M{
		"uuid": id,
		"t": bson.M{
			"$gte": from.Unix(),
			"$lte": to.Unix(),
		},
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find())
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
