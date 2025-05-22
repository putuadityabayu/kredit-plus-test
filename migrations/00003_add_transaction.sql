-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transactions (
	id UUID NOT NULL PRIMARY KEY,
	contract_number VARCHAR(255) NOT NULL COMMENT 'Nomor kontrak (unik)',
	user_id UUID NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
	otr DECIMAL NOT NULL COMMENT 'Nominal on the road',
	admin_fee DECIMAL NOT NULL COMMENT 'Admin fee',
	installment_amount DECIMAL NOT NULL COMMENT 'Jumlah Cicilan per bulan',
	interest_amount DECIMAL NOT NULL COMMENT 'Total Bunga yang ditagihkan',
	asset_name VARCHAR(255) NOT NULL COMMENT 'Nama Asset yang dibeli',
	tenor INT NOT NULL COMMENT 'Tenor dalam bulan',
	transaction_date TIMESTAMP NOT NULL COMMENT 'Tanggal transaksi',
	status ENUM('pending', 'approved', 'rejected') NOT NULL DEFAULT 'pending',
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

	INDEX idx_transactions (transaction_date desc),
	UNIQUE KEY unique_contract (contract_number)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd
