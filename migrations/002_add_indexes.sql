-- Migration: Create indexes for better query performance

CREATE INDEX idx_transactions_customer_timestamp ON transactions(customer_id, timestamp DESC);
CREATE INDEX idx_users_token_expiry ON users(token, token_expiry);
CREATE INDEX idx_customers_created_at ON customers(created_at DESC);
CREATE INDEX idx_transactions_created_at ON transactions(created_at DESC);
