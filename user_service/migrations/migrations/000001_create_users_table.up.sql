CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    balance NUMERIC(10, 2) DEFAULT 0,
    role TEXT DEFAULT 'user'
);
