-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `users` (
	id UUID NOT NULL PRIMARY KEY,
	nik VARCHAR(16),
	full_name VARCHAR(255) NOT NULL,
	legal_name VARCHAR(255),
	birth_place VARCHAR(255),
	birth_date DATE,
	salary DECIMAL,
	ktp_photo_url VARCHAR(255),
	selfie_photo_url VARCHAR(255),
	password VARCHAR(255) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP NULL,

	INDEX idx_users (deleted_at, created_at),
	UNIQUE INDEX idx_users_nik (nik)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users
-- +goose StatementEnd
