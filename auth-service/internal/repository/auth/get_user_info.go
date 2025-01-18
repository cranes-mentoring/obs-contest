package auth

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	models "github.com/cranes-mentoring/obs-contest/auth-service/internal/model"
)

const query = `
		SELECT username, login, ip_address, email, device, country, name
		FROM users
		WHERE username = $1
	`

// GetUserInfo fetches user information from the database.
func (r *UserRepository) GetUserInfo(ctx context.Context, username string) (*models.User, error) {
	tracer := otel.Tracer("auth-service")
	ctx, span := tracer.Start(ctx, "UserRepository.GetUserInfo", trace.WithAttributes(
		attribute.String("db.statement", query),
		attribute.String("db.user", username),
	))
	defer span.End()

	var user models.User
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.Username,
		&user.Login,
		&user.IPAddress,
		&user.Email,
		&user.Device,
		&user.Country,
		&user.Name,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to execute query")
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}

	span.SetStatus(codes.Ok, "Query executed successfully")
	return &user, nil
}
