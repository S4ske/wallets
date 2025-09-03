CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    balance INTEGER NOT NULL CHECK (balance >= 0) DEFAULT 0
);