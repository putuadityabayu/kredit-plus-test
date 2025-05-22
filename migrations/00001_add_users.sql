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

	INDEX idx_users (deleted_at, created_at desc),
	UNIQUE INDEX idx_users_nik (nik)
);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT IGNORE INTO users (
	id, nik, full_name, legal_name, birth_place, birth_date,
	salary, ktp_photo_url, selfie_photo_url, password,
	created_at, updated_at, deleted_at
)
VALUES
	('0196f7ef-49de-79e9-b5a6-227b15de5240', '1234567890112345', 'Budi', 'Budi Legal Name', 'Mataram', '2000-01-01', 7500000, NULL, NULL, '$2a$08$nuEOGocBbWf6e3NNR4hBt.ben3okjOXFZUjLmXN3TuE8HTfrhyBeu', '2025-05-22 12:19:36', '2025-05-22 12:19:36', NULL),
	('0196f7ef-6f01-7685-b2f4-b5a796c2ce15', '1234567890111234', 'Annisa', 'Annisa Legal Name', 'Malang', '1995-01-01', 9000000, NULL, NULL, '$2a$08$tVSs0eqQNic7LULEHeBRjOnqEw0jy86T6Rz7KWk6b/CAgvfX2Pw/S', '2025-05-22 12:19:46', '2025-05-22 12:19:46', NULL);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users
-- +goose StatementEnd
