package main

import (
	"MessageService/internal/domains/message/handler/kafka"
	"context"
	"log"

	msgHttp "MessageService/internal/domains/message/handler/http"
	"MessageService/internal/domains/message/repository/mongodb"
	"MessageService/internal/domains/message/usecase"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	// Подключение к Mongo
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("messages")

	// Репозиторий и UseCase
	repo := mongodb.NewMessageRepository(db, "messages")
	uc := usecase.NewMessageUseCase(repo)

	// Gin router
	router := gin.Default()

	// Регистрируем хендлер
	msgHandler := msgHttp.NewMessageHandler(uc)
	msgHandler.RegisterRoutes(router)
	kafkaConsumer := kafka.NewMessageConsumer(
		[]string{"localhost:9092"},
		"test-topic",
		"test-group",
		uc,
	)

	go func() {
		if err := kafkaConsumer.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("HTTP server started on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}

	// Блокируем main
	select {}
}
