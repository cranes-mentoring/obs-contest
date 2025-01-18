package auth

import (
	"context"
	"fmt"

	models "github.com/cranes-mentoring/obs-contest/auth-service/internal/model"
)

const query = `
		SELECT username, login, ip_address, email, device, country, name
		FROM users
		WHERE username = $1
	`

// GetUserInfo fetches user information from the database.
func (r *UserRepository) GetUserInfo(ctx context.Context, username string) (*models.User, error) {
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
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}

	return &user, nil
}
