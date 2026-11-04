CREATE TABLE transactions
(
    `uuid` UUID,
    `bank_id` String,
    `sender_account` String,
    `receiver_account` String,
    `amount` Decimal(18, 2),
    `currency` LowCardinality(String),
    `timestamp` DateTime64(6, 'UTC'),
    `update_at` DateTime64(6, 'UTC'),
    `status` LowCardinality(String)
)
ENGINE = ReplacingMergeTree(update_at)
ORDER BY (bank_id, timestamp, uuid);