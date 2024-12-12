package service

import (
	"context"
	"time"

	"github.com/Alarmtekgit/websocket/internal/domain"
	"github.com/Alarmtekgit/websocket/internal/repository"
)

type HistoryService interface {
	CreateHistory(ctx context.Context, history *domain.History) error
	FindHistoryByID(ctx context.Context, equipmentid string) (*domain.History, error)
	GetHistoryByID(ctx context.Context, messageType string) ([]*domain.History, error)
}

type historyService struct {
	repo repository.HistoryRepository
}

func NewHistoryService(repo repository.HistoryRepository) HistoryService {
	return &historyService{repo: repo}
}

func (s *historyService) CreateHistory(ctx context.Context, history *domain.History) error {
	history.Timestamp = time.Now()
	return s.repo.CreateHistory(ctx, history)
}

func (s *historyService) FindHistoryByID(ctx context.Context, equipmentid string) (*domain.History, error) {
	return s.repo.FindHistoryByID(ctx, equipmentid)
}

func (s *historyService) GetHistoryByID(ctx context.Context, equipmentid string) ([]*domain.History, error) {
	return s.repo.GetHistoryByID(ctx, equipmentid)
}
