package main

import (
	"context"
	"net"
	"os"
	"os/signal"

	pb "github.com/cranes-mentoring/obs-contest/auth-service/generated/auth-service/proto"
	"github.com/cranes-mentoring/obs-contest/auth-service/internal/db"
	"github.com/cranes-mentoring/obs-contest/auth-service/internal/logging"
	auth_repo "github.com/cranes-mentoring/obs-contest/auth-service/internal/repository/auth"
	"github.com/cranes-mentoring/obs-contest/auth-service/internal/service/auth"
	"github.com/cranes-mentoring/obs-contest/auth-service/internal/tracing"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	_ "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const connStr = "postgres://postgres:postgres@postgres:5432/public"

func main() {
	ctx := context.Background()

	logging.SetupLogger()
	logger := logging.Logger
	defer logger.Sync()

	shutdown := tracing.InitTracer(ctx)
	defer shutdown()

	dbpool := db.InitDB(ctx, connStr)
	defer dbpool.Close()

	authRepo := auth_repo.NewUserRepository(dbpool)
	authService := auth.NewUserService(authRepo, logger)

	server := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))

	pb.RegisterAuthServiceServer(server, authService)

	reflection.Register(server)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		logger.Info("Starting gRPC server on port 50051...")
		if err := server.Serve(listener); err != nil {
			logger.Fatal("failed to serve", zap.Error(err))
		}
	}()

	<-stop
	logger.Info("Shutting down server...")
	server.GracefulStop()
}
