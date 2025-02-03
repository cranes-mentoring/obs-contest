package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/cranes-mentoring/obs-contest/purchase-processor/internal/logging"
	"github.com/cranes-mentoring/obs-contest/purchase-processor/internal/service"
	"github.com/cranes-mentoring/obs-contest/purchase-processor/internal/tracing"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
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

	conn, err := grpc.NewClient(
		"auth-service:50051",
		grpc.WithInsecure(),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		logger.Fatal("error", zap.Error(err))
	}

	defer conn.Close()

	// init layers
	authServicePb := pb.NewAuthServiceClient(conn)
	authService := service.NewUserService(authServicePb)
	handler := purchase_handler.NewConsumerGroupHandler(*authService)

	server := &http.Server{
		Addr: ":8084",
	}

	go func() {
		logger.Info("Server running.")
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("HTTP server ListenAndServe error", zap.Error(err))
		}
	}()

	brokers := strings.Split(getEnv("KAFKA_BROKERS", "kafka:9092"), ",")
	topic := getEnv("KAFKA_TOPIC", "purchases")
	group := getEnv("KAFKA_GROUP", "purchase-processor-group")

	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		logger.Fatal("Error creating consumer group", zap.Error(err))
	}
	defer func() {
		if err := consumerGroup.Close(); err != nil {
			logger.Fatal("Error closing consumer group", zap.Error(err))
		}
	}()

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		for {
			if err := consumerGroup.Consume(ctx, []string{topic}, handler); err != nil {
				logger.Fatal("Error from consumer", zap.Error(err))
			}
		}
	}()

	logger.Info("Started consumer for topic", zap.String("topic", topic))

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
