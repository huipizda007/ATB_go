-- schema.sql
CREATE TABLE fruits (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    brand       TEXT, 
    price_per_kg NUMERIC(10, 2) NOT NULL,
    stock_kg    INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);