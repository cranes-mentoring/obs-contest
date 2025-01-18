-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(50) NOT NULL,
                       login VARCHAR(50) NOT NULL,
                       ip_address VARCHAR(50) NOT NULL,
                       email VARCHAR(100) NOT NULL,
                       device VARCHAR(50),
                       country VARCHAR(50),
                       name VARCHAR(100)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
