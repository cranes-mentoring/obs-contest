package service

import (
	"context"
	"errors"

	pb "github.com/cranes-mentoring/obs-contest/purchase-processor/generated/auth-service/proto"
)

type UserService struct {
	authService pb.AuthServiceClient
}

func NewUserService(authService pb.AuthServiceClient) *UserService {
	return &UserService{authService: authService}
}

func (s *UserService) FindUser(ctx context.Context, username string) (string, error) {
	if username == "" {
		return "", errors.New("empty string")
	}

	info, err := s.authService.GetUserInfo(ctx, &pb.GetUserInfoRequest{
		Username: username,
	})

	if err != nil {
		return "", err
	}

	return info.Email, err
}
