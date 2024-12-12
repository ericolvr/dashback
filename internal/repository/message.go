package repository

import (
	"context"
	"errors"

	"github.com/Alarmtekgit/websocket/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrNotFound = errors.New("message not found")

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *domain.Message) error
	FindMessageByID(ctx context.Context, id primitive.ObjectID) (*domain.Message, error)
	FindByFields(ctx context.Context, filter bson.M) (*domain.Message, error)
	UpdateMessage(ctx context.Context, id primitive.ObjectID, monitored bool) error
	DeleteMessage(ctx context.Context, id primitive.ObjectID) error
	GetMessagesByType(ctx context.Context, messageType string) ([]*domain.Message, error)
	GetMonitoredMessages(ctx context.Context, monitored bool) ([]*domain.Message, error)
}

type messageRepository struct {
	collection *mongo.Collection
}

func NewMessageRepository(db *mongo.Database) MessageRepository {
	return &messageRepository{
		collection: db.Collection("messages"),
	}
}

func (r *messageRepository) CreateMessage(ctx context.Context, message *domain.Message) error {
	_, err := r.collection.InsertOne(ctx, message)
	return err
}

func (r *messageRepository) FindMessageByID(ctx context.Context, id primitive.ObjectID) (*domain.Message, error) {
	var message domain.Message
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&message)
	return &message, err
}

func (r *messageRepository) FindByFields(ctx context.Context, filter bson.M) (*domain.Message, error) {
	var message domain.Message
	err := r.collection.FindOne(ctx, filter).Decode(&message)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("message not found")
		}
		return nil, err
	}
	return &message, nil
}

func (r *messageRepository) UpdateMessage(ctx context.Context, id primitive.ObjectID, monitored bool) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"monitored": monitored}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *messageRepository) DeleteMessage(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *messageRepository) GetMessagesByType(ctx context.Context, messageType string) ([]*domain.Message, error) {
	filter := bson.M{
		"type":      messageType,
		"monitored": false,
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*domain.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *messageRepository) GetMonitoredMessages(ctx context.Context, monitored bool) ([]*domain.Message, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"monitored": monitored})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*domain.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}
