package usecase

import (
	"context"
	"time"

	"MessageService/internal/domains/message/model"
	"MessageService/internal/domains/message/repository/mongodb"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageUseCase interface {
	Save(ctx context.Context, msg *model.Message) error
	GetByID(ctx context.Context, id string) (*model.Message, error)
	GetByTimeRange(ctx context.Context, from, to time.Time) ([]*model.Message, error)
}

type messageUseCase struct {
	repo mongodb.MessageRepository
}

func NewMessageUseCase(repo mongodb.MessageRepository) MessageUseCase {
	return &messageUseCase{
		repo: repo,
	}
}

func (uc *messageUseCase) Save(ctx context.Context, msg *model.Message) error {
	return uc.repo.Insert(ctx, msg)
}

func (uc *messageUseCase) GetByID(ctx context.Context, id string) (*model.Message, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return uc.repo.GetByID(ctx, objID)
}

func (uc *messageUseCase) GetByTimeRange(ctx context.Context, from, to time.Time) ([]*model.Message, error) {
	return uc.repo.GetByTimeRange(ctx, from, to)
}
