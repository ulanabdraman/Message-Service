// internal/domains/message/usecase/usecase_test.go
package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"MessageService/internal/domains/message/model"
	"MessageService/internal/domains/message/usecase"
)

// --- мок репозитория ---
type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Insert(ctx context.Context, msg *model.Message) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *mockRepo) GetByID(ctx context.Context, id int64) (*model.Message, error) {
	return nil, nil
}

func (m *mockRepo) GetByTimeRange(ctx context.Context, id int64, from time.Time, to time.Time) ([]*model.Message, error) {
	return nil, nil
}

// --- тест ---
func TestSave_CallsRepo(t *testing.T) {
	mrepo := new(mockRepo)
	uc := usecase.NewMessageUseCase(mrepo)

	msg := &model.Message{DT: 123, ST: 1}
	mrepo.On("Insert", mock.Anything, msg).Return(nil)

	err := uc.Save(context.Background(), msg)

	assert.NoError(t, err)
	mrepo.AssertCalled(t, "Insert", mock.Anything, msg)
}
