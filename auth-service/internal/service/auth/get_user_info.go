package auth

import (
	"context"

	pb "github.com/cranes-mentoring/obs-contest/auth-service/generated/auth-service/proto"
	"github.com/cranes-mentoring/obs-contest/auth-service/internal/logging"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errNotFound = "error fetching user info: %v"

func (s *UserService) GetUserInfo(ctx context.Context, request *pb.GetUserInfoRequest) (*pb.UserInfoResponse, error) {
	username := request.GetUsername()

	tracedLogger := logging.AddTraceContextToLogger(ctx)

	tracedLogger.Info("find user info",
		zap.String("operation", "find user"),
		zap.String("username", username),
	)

	user, err := s.userRepo.GetUserInfo(ctx, username)
	if err != nil {
		s.logger.Error(errNotFound)

		return nil, status.Errorf(codes.NotFound, errNotFound, err)
	}

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
