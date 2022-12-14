CREATE TABLE IF NOT EXISTS payloads (
    id_payload SERIAL,
    user_id INTEGER NOT NULL,
    order_id VARCHAR(256),
    transaction_time VARCHAR(256),
    transaction_status VARCHAR(256),
    transaction_id VARCHAR(256),
    status_code VARCHAR(256),
    signature_key VARCHAR(256),
    settlement_time VARCHAR(256),
    payment_type VARCHAR(256),
    merchant_id VARCHAR(256),
    gross_amount VARCHAR(256),
    fraud_status VARCHAR(256),
    bank_type VARCHAR(256),
    va_number VARCHAR(256),
    biller_code VARCHAR(256),
    bill_key VARCHAR(256),
    receipt_number VARCHAR(256),
    address VARCHAR(256),
    courier VARCHAR(30),
    courier_service VARCHAR(100),
    PRIMARY KEY (id_payload),
    FOREIGN KEY (user_id) REFERENCES users(id_user) ON DELETE CASCADE ON UPDATE CASCADE
)