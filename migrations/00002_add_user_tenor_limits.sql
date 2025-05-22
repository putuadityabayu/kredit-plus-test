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

-- +goose StatementBegin
INSERT IGNORE INTO user_tenor_limits (id, user_id, tenor_in_months, limit_amount, created_at, updated_at)
VALUES  ('905b12ff-b4e7-45a3-8864-1a5514530625', '0196f7ef-49de-79e9-b5a6-227b15de5240', 6, 700000, '2025-05-22 19:25:15', '2025-05-22 19:25:15'),
		('426aadf7-7845-4783-87b5-2153eede923a', '0196f7ef-49de-79e9-b5a6-227b15de5240', 3, 500000, '2025-05-22 19:25:15', '2025-05-22 19:25:15'),
		('4261f7fe-c84c-4cdd-bf66-51838fd0ec0a', '0196f7ef-6f01-7685-b2f4-b5a796c2ce15', 1, 1000000, '2025-05-22 19:27:05', '2025-05-22 19:27:05'),
		('f564df9b-ccf4-4ce8-8442-59768b024c7a', '0196f7ef-6f01-7685-b2f4-b5a796c2ce15', 2, 1200000, '2025-05-22 19:27:05', '2025-05-22 19:27:05'),
		('c1688126-1a9e-4649-8482-635bd5c6c8e6', '0196f7ef-6f01-7685-b2f4-b5a796c2ce15', 3, 1500000, '2025-05-22 19:27:05', '2025-05-22 19:27:05'),
		('fa436cc8-7539-4f57-991b-890637125bc7', '0196f7ef-49de-79e9-b5a6-227b15de5240', 2, 200000, '2025-05-22 19:25:15', '2025-05-22 19:25:15'),
		('3f2fbacb-e669-4382-8e64-933e15c4503f', '0196f7ef-6f01-7685-b2f4-b5a796c2ce15', 6, 2000000, '2025-05-22 19:27:05', '2025-05-22 19:27:05'),
		('de0f0668-7a72-4ac8-befb-b3d3de8ac857', '0196f7ef-49de-79e9-b5a6-227b15de5240', 1, 100000, '2025-05-22 19:25:15', '2025-05-22 19:25:15');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_tenor_limits;
-- +goose StatementEnd
