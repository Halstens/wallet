CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),  -- Автогенерация
    balance BIGINT NOT NULL DEFAULT 0
);

INSERT INTO wallets (balance) VALUES (1000);