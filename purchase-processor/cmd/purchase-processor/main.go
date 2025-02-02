package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/cranes-mentoring/obs-contest/purchase-processor/internal/logging"
	"github.com/cranes-mentoring/obs-contest/purchase-processor/internal/middleware"
	"github.com/cranes-mentoring/obs-contest/purchase-processor/internal/service"
	"github.com/cranes-mentoring/obs-contest/purchase-processor/internal/tracing"
	"google.golang.org/grpc"

	pb "github.com/cranes-mentoring/obs-contest/purchase-processor/generated/auth-service/proto"
	purchase_handler "github.com/cranes-mentoring/obs-contest/purchase-processor/internal/handler"
)

func main() {
	ctx := context.Background()

	logging.SetupLogger()
	logger := logging.Logger
	defer logger.Sync()

	shutdown := tracing.InitTracer(ctx)
	defer shutdown()

	brokers := strings.Split(getEnv("KAFKA_BROKERS", "0.0.0.0:9092"), ",")
	topic := getEnv("KAFKA_TOPIC", "purchases")
	group := getEnv("KAFKA_GROUP", "purchase-processor-group")

	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	defer func() {
		if err := consumerGroup.Close(); err != nil {
			log.Fatalf("Error closing consumer group: %v", err)
		}
	}()

	conn, err := grpc.NewClient(
		"auth-service:50051",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(middleware.ClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer conn.Close()

	// init layers
	authServicePb := pb.NewAuthServiceClient(conn)
	authService := service.NewUserService(authServicePb)
	handler := purchase_handler.NewConsumerGroupHandler(*authService)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		for {
			if err := consumerGroup.Consume(ctx, []string{topic}, handler); err != nil {
				log.Fatalf("Error from consumer: %v", err)
			}
		}
	}()

	log.Printf("Started consumer for topic: %s", topic)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	cancel()

	log.Println("Shutting down consumer...")
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
