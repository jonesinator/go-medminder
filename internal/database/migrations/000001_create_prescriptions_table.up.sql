CREATE TABLE IF NOT EXISTS prescriptions (
    name VARCHAR(64) PRIMARY KEY,
    quantity DECIMAL(10, 2),
    rate DECIMAL(10, 2),
    updated DATE
);