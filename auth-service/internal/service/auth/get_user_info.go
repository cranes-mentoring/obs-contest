package auth

import (
	"context"

	pb "github.com/cranes-mentoring/obs-contest/auth-service/generated/auth-service/proto"
	"github.com/cranes-mentoring/obs-contest/auth-service/internal/logging"

	"go.opentelemetry.io/otel"
	otelCodes "go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errNotFound = "error fetching user info: %v"

func (s *UserService) GetUserInfo(ctx context.Context, request *pb.GetUserInfoRequest) (*pb.UserInfoResponse, error) {
	username := request.GetUsername()

	tracer := otel.Tracer("auth-service")
	ctx, span := tracer.Start(ctx, "GetUserInfo")
	defer span.End()

	tracedLogger := logging.AddTraceContextToLogger(ctx)

	tracedLogger.Info("find user info",
		zap.String("operation", "find user"),
		zap.String("username", username),
	)

	user, err := s.userRepo.GetUserInfo(ctx, username)
	if err != nil {
		s.logger.Error(errNotFound)
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, "Failed to fetch user info")

		return nil, status.Errorf(grpcCodes.NotFound, errNotFound, err)
	}

	span.SetStatus(otelCodes.Ok, "User info retrieved successfully")

	return &pb.UserInfoResponse{
		Username:  user.Username,
		Login:     user.Login,
		IpAddress: user.IPAddress,
		Email:     user.Email,
		Device:    user.Device,
		Country:   user.Country,
		Name:      user.Name,
	}, nil
}
