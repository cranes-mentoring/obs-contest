package service

import (
	"context"
	"errors"

	pb "github.com/cranes-mentoring/obs-contest/purchase-processor/generated/auth-service/proto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type UserService struct {
	authService pb.AuthServiceClient
}

func NewUserService(authService pb.AuthServiceClient) *UserService {
	return &UserService{authService: authService}
}

func (s *UserService) findUser(ctx context.Context, username string) (string, error) {
	if username == "" {
		return "", errors.New("empty string")
	}

	ctx, span := s.handleTracing(ctx, username)
	defer span.End()

	info, err := s.authService.GetUserInfo(ctx, &pb.GetUserInfoRequest{
		Username: username,
	})

	if err != nil {
		return "", err
	}

	span.SetStatus(codes.Ok, "Auth called successfully")

	return info.Email, err
}

func (s *UserService) handleTracing(ctx context.Context, username string) (context.Context, trace.Span) {
	tracer := otel.Tracer("purchase-processor")

	ctx, span := tracer.Start(ctx, "ProcessProcessor.findUser", trace.WithAttributes(
		attribute.String("username", username),
	))

	return ctx, span
}
