CREATE TABLE IF NOT EXISTS bets (
    id UUID PRIMARY KEY,
    user_id TEXT NOT NULL,
    event_id TEXT NOT NULL,
    amount NUMERIC(10, 2),
    odds NUMERIC(5, 2),
    status TEXT,
    payout NUMERIC(10, 2),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);
