CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    start_time TIMESTAMPTZ,
    status TEXT,
    winner_id TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);
