CREATE TABLE transactions(
	uuid UUID PRIMARY KEY,
	bank_id VARCHAR(10),
	sender_account VARCHAR(20),
	receiver_account VARCHAR(20),
	amount DECIMAL(18, 2),
	currency VARCHAR(5),
	occurred_at TIMESTAMP WITH TIME ZONE,
	status VARCHAR(10)
)