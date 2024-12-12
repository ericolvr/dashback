package service

import (
	"context"
	"time"

	"github.com/Alarmtekgit/websocket/internal/domain"
	"github.com/Alarmtekgit/websocket/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageService interface {
	CreateMessage(ctx context.Context, message *domain.Message) error
	InsertIfNotExists(ctx context.Context, message *domain.Message) error
	FindMessageByID(ctx context.Context, id primitive.ObjectID) (*domain.Message, error)
	UpdateMessage(ctx context.Context, id primitive.ObjectID, monitored bool) error
	DeleteMessage(ctx context.Context, id primitive.ObjectID) error
	GetMessagesByType(ctx context.Context, messageType string) ([]*domain.Message, error)
	GetMonitoredMessages(ctx context.Context, monitored bool) ([]*domain.Message, error)
}

type messageService struct {
	repo repository.MessageRepository
}

func NewMessageService(repo repository.MessageRepository) MessageService {
	return &messageService{repo: repo}
}

func (s *messageService) CreateMessage(ctx context.Context, message *domain.Message) error {
	message.Timestamp = time.Now()
	return s.repo.CreateMessage(ctx, message)
}

func (s *messageService) InsertIfNotExists(ctx context.Context, message *domain.Message) error {
	filter := bson.M{
		"uniorg": message.Uniorg,
		"panel":  message.Panel,
		"node":   message.Node,
		"type":   message.Type,
		"status": message.Status,
	}

	existingMessage, err := s.repo.FindByFields(ctx, filter)

	if err == nil && existingMessage != nil {
		return nil
	}

	return s.repo.CreateMessage(ctx, message)
}

func (s *messageService) FindMessageByID(ctx context.Context, id primitive.ObjectID) (*domain.Message, error) {
	return s.repo.FindMessageByID(ctx, id)
}

func (s *messageService) UpdateMessage(ctx context.Context, id primitive.ObjectID, monitored bool) error {
	return s.repo.UpdateMessage(ctx, id, monitored)
}

func (s *messageService) DeleteMessage(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.DeleteMessage(ctx, id)
}

func (s *messageService) GetMessagesByType(ctx context.Context, messageType string) ([]*domain.Message, error) {
	return s.repo.GetMessagesByType(ctx, messageType)
}

func (s *messageService) GetMonitoredMessages(ctx context.Context, monitored bool) ([]*domain.Message, error) {
	return s.repo.GetMonitoredMessages(ctx, monitored)
}
