-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_tenor_limits (
	id UUID NOT NULL PRIMARY KEY,
	user_id UUID NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
	tenor_in_months INT NOT NULL COMMENT 'tenor dalam bulan (e.g., 1,2,3,6)',
	limit_amount DECIMAL NOT NULL COMMENT 'Jumlah limit untuk tenor ini',
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_tenor_limits;
-- +goose StatementEnd
