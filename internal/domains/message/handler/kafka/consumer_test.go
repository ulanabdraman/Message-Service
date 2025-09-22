package kafka

import (
	"context"
	"testing"
	"time"

	"MessageService/internal/domains/message/model"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---- мок на usecase.MessageUseCase ----
type mockUseCase struct{ mock.Mock }

func (m *mockUseCase) Save(ctx context.Context, msg *model.Message) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}
func (m *mockUseCase) GetByID(ctx context.Context, id int64) (*model.Message, error) {
	args := m.Called(ctx, id)
	if v := args.Get(0); v != nil {
		return v.(*model.Message), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockUseCase) GetByTimeRange(ctx context.Context, from, to time.Time, id int64) ([]*model.Message, error) {
	args := m.Called(ctx, from, to)
	if v := args.Get(0); v != nil {
		return v.([]*model.Message), args.Error(1)
	}
	return nil, args.Error(1)
}

// ---- 1) проверяем корректность dtoToModel ----
func Test_dtoToModel_OK(t *testing.T) {
	dto := MessageDTO{
		T:  1712345678,
		ST: 7,
		Pos: PosDTO{
			X: 51.15, Y: 71.42, Z: 123, A: 180, S: 10, St: 3,
		},
		Params: map[string]interface{}{"foo": "bar"},
	}

	got := dtoToModel(dto)

	want := &model.Message{
		T:  1712345678,
		ST: 7,
		Pos: model.Pos{
			X: 51.15, Y: 71.42, Z: 123, A: 180, S: 10, St: 3,
		},
		Params: map[string]interface{}{"foo": "bar"},
	}

	assert.Equal(t, want, got)
}

// ---- 2) проверяем, что NewMessageConsumer правильно кладёт useCase внутрь ----
func Test_NewMessageConsumer_Injection(t *testing.T) {
	muc := new(mockUseCase)
	c := NewMessageConsumer([]string{"localhost:9092"}, "topic", "group", muc)

	// поле называется useCase (с маленькой буквы) и неэкспортируемое,
	// но мы в том же пакете, поэтому видим его
	assert.NotNil(t, c)
	assert.Equal(t, muc, c.useCase)

	// на всякий случай проверим, что reader тоже создался
	assert.IsType(t, &kafka.Reader{}, c.reader)
	_ = c.Close()
}
