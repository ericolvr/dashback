package repository

import (
	"context"

	"github.com/Alarmtekgit/websocket/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type HistoryRepository interface {
	CreateHistory(ctx context.Context, history *domain.History) error
	FindHistoryByID(ctx context.Context, equipmentid string) (*domain.History, error)
	GetHistoryByID(ctx context.Context, equipmentid string) ([]*domain.History, error)
}

type historyRepository struct {
	collection *mongo.Collection
}

func NewHistoryRepository(db *mongo.Database) HistoryRepository {
	return &historyRepository{
		collection: db.Collection("history"),
	}
}

func (r *historyRepository) CreateHistory(ctx context.Context, history *domain.History) error {
	_, err := r.collection.InsertOne(ctx, history)
	return err
}

func (r *historyRepository) FindHistoryByID(ctx context.Context, equipmentid string) (*domain.History, error) {
	var history domain.History
	err := r.collection.FindOne(ctx, bson.M{"equipmentid": equipmentid}).Decode(&history)
	return &history, err
}

func (r *historyRepository) GetHistoryByID(ctx context.Context, equipmentid string) ([]*domain.History, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"equipmentid": equipmentid})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var histories []*domain.History
	if err := cursor.All(ctx, &histories); err != nil {
		return nil, err
	}
	return histories, nil
}
