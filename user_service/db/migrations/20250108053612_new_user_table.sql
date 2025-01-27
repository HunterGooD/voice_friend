-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY, 
    login VARCHAR(160) NOT NULL UNIQUE, 
    name VARCHAR(160) NOT NULL, 
    uid UUID DEFAULT uuid_generate_v4(), 
    email VARCHAR(255) NOT NULL UNIQUE, 
    password VARCHAR(255) NOT NULL, 
    is_active BOOLEAN DEFAULT TRUE, 
    last_login TIMESTAMP DEFAULT NULL, 
    role VARCHAR(20) DEFAULT 'user', 
    profile_picture VARCHAR(255), 
    phone VARCHAR(20),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX idx_users_uid ON users (uid);
CREATE INDEX idx_users_is_active ON users (is_active);
CREATE INDEX idx_users_is_delete ON users (deleted_at);
-- CREATE INDEX idx_users_last_login ON users (last_login DESC);
-- CREATE INDEX idx_users_phone ON users (phone);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
