CREATE TABLE IF NOT EXISTS monitors (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    target_url TEXT NOT NULL,
    interval_seconds INT NOT NULL,
    failure_threshold INT NOT NULL,
    consecutive_failures INT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL,
    last_checked_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS incidents (
    id VARCHAR(36) PRIMARY KEY,
    monitor_id VARCHAR(36) NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    status VARCHAR(20) NOT NULL,
    opened_at TIMESTAMPTZ NOT NULL,
    resolved_at TIMESTAMPTZ NULL
);