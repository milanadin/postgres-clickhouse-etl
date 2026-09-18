CREATE TABLE transactions(
	uuid UUID PRIMARY KEY,
	bank_id VARCHAR(10),
	from_account VARCHAR(20),
	to_account VARCHAR(20),
	amount DECIMAL(18, 2),
	currency INT,
	unixtime BIGINT,
	status INT
)