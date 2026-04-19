CREATE TABLE IF NOT EXISTS project_profiles (
    content_id BIGINT PRIMARY KEY REFERENCES content_items(id) ON DELETE CASCADE,
    project_name VARCHAR(255) NOT NULL DEFAULT '',
    stack TEXT NOT NULL DEFAULT '',
    repo_url TEXT NOT NULL DEFAULT '',
    demo_url TEXT NOT NULL DEFAULT '',
    started_at DATE NULL,
    ended_at DATE NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS experience_profiles (
    content_id BIGINT PRIMARY KEY REFERENCES content_items(id) ON DELETE CASCADE,
    org_name VARCHAR(255) NOT NULL DEFAULT '',
    role_name VARCHAR(255) NOT NULL DEFAULT '',
    location VARCHAR(255) NOT NULL DEFAULT '',
    started_at DATE NULL,
    ended_at DATE NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS timeline_event_profiles (
    content_id BIGINT PRIMARY KEY REFERENCES content_items(id) ON DELETE CASCADE,
    event_time TIMESTAMPTZ NULL,
    event_category VARCHAR(64) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS insight_profiles (
    content_id BIGINT PRIMARY KEY REFERENCES content_items(id) ON DELETE CASCADE,
    insight_level VARCHAR(32) NOT NULL DEFAULT 'general',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS portfolio_profiles (
    content_id BIGINT PRIMARY KEY REFERENCES content_items(id) ON DELETE CASCADE,
    artifact_type VARCHAR(64) NOT NULL DEFAULT '',
    external_link TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
