-- +goose Up
-- +goose StatementBegin
INSERT INTO users (username, login, ip_address, email, device, country, name)
VALUES
    ('first_user', 'first_login', '127.0.0.1', 'user@example.com', 'PC', 'USA', 'John Doe');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE id = 1;
-- +goose StatementEnd
