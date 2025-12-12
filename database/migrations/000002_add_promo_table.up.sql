-- file: migrations/xxxxxx_add_promo_table.up.sql

CREATE TABLE promos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    discount_type VARCHAR(20) NOT NULL,
    discount_value DECIMAL(10, 2) NOT NULL,
    valid_until TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

ALTER TABLE transactions 
ADD COLUMN promo_id UUID REFERENCES promos(id),
ADD COLUMN discount_amount DECIMAL(10, 2) DEFAULT 0,
ADD COLUMN final_amount DECIMAL(10, 2) DEFAULT 0;