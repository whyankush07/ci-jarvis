CREATE TABLE IF NOT EXISTS runs (
    id uuid PRIMARY KEY UNIQUE,
    repo_url TEXT NOT NULL,
    pr_url TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    current_step TEXT,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_runs_status ON runs(status);