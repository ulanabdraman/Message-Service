package mongodb_test

import (
	"context"
	"log"
	"testing"
	"time"

	"MessageService/internal/domains/message/model"
	"MessageService/internal/domains/message/repository/mongodb"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestInsertAndGetByID(t *testing.T) {
	ctx := context.Background()

	// Подключение к MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	db := client.Database("testdb")
	repo := mongodb.NewMessageRepository(db, "messages")

	// Тестируемое сообщение
	msg := &model.Message{
		UUID: 3453,
		T:    time.Now().Unix(),
		ST:   1,
		Pos:  model.Pos{X: 1.1, Y: 2.2, Z: 3, A: 4, S: 5, St: 6},
		Params: map[string]interface{}{
			"custom": "value",
			"speed":  42,
		},
	}
	log.Println(msg.T)

	// Вставка
	err = repo.Insert(ctx, msg)
	require.NoError(t, err)

	// Получение по ID
	fetched, err := repo.GetByID(ctx, msg.UUID)
	require.NoError(t, err)
	require.NotNil(t, fetched)

	// Проверка полей (упрощённо)
	require.Equal(t, msg.UUID, fetched.UUID)
	require.Equal(t, msg.T, fetched.T)
	require.Equal(t, msg.ST, fetched.ST)
	require.Equal(t, msg.Pos, fetched.Pos)

	// Очистка тестовых данных
	_, _ = db.Collection("messages").DeleteOne(ctx, primitive.M{"uuid": msg.UUID})
}

func TestGetByTimeRange(t *testing.T) {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	db := client.Database("testdb")
	repo := mongodb.NewMessageRepository(db, "messages")

	// используем уникальный uuid для теста
	testUUID := time.Now().UnixNano()

	now := time.Now().UTC().Truncate(time.Second)

	msgs := []*model.Message{
		{UUID: testUUID, T: now.Add(-2 * time.Minute).Unix(), ST: 1, Pos: model.Pos{X: 1, Y: 1}},
		{UUID: testUUID, T: now.Unix(), ST: 2, Pos: model.Pos{X: 2, Y: 2}}, // должно попасть
		{UUID: testUUID, T: now.Add(2 * time.Minute).Unix(), ST: 3, Pos: model.Pos{X: 3, Y: 3}},
	}

	for _, m := range msgs {
		err := repo.Insert(ctx, m)
		require.NoError(t, err)
	}

	from := now.Add(-1 * time.Minute)
	to := now.Add(1 * time.Minute)

	results, err := repo.GetByTimeRange(ctx, testUUID, from, to)
	require.NoError(t, err)

	require.Len(t, results, 1, "должно вернуться только одно сообщение")
	require.Equal(t, msgs[1].T, results[0].T)

	// очистка
	_, _ = db.Collection("messages").DeleteMany(ctx, primitive.M{"uuid": testUUID})
}
