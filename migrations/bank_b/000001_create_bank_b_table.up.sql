CREATE TABLE transactions(
	uuid UUID PRIMARY KEY,
	bank_id VARCHAR(10),
	from_account VARCHAR(20),
	to_account VARCHAR(20),
	amount BIGINT,
	currency VARCHAR(5),
	transaction_time TIMESTAMP WITH TIME ZONE,
	status VARCHAR(20)
)