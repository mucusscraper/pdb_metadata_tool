-- +goose Up
CREATE TABLE entries (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    rcsb_id TEXT NOT NULL UNIQUE,
    deposit_date TEXT NOT NULL DEFAULT 'Unknown',
    doi TEXT NOT NULL DEFAULT 'Unknown',
    paper_title TEXT NOT NULL DEFAULT 'Unknown',
    method TEXT NOT NULL DEFAULT 'Unknown',
    user_group TEXT NOT NULL DEFAULT ''
);

-- +goose Down
DROP TABLE entries;
