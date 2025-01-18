package auth

import (
	"github.com/cranes-mentoring/obs-contest/auth-service/generated/auth-service/proto"
	"github.com/cranes-mentoring/obs-contest/auth-service/internal/repository/auth"
	"go.uber.org/zap"
)

type UserService struct {
	proto.UnimplementedAuthServiceServer

	userRepo *auth.UserRepository
	logger   *zap.Logger
}

// NewUserService creates a new UserService.
func NewUserService(userRepo *auth.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{userRepo: userRepo, logger: logger}
}
