CREATE TABLE IF NOT EXISTS todos (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL,
    title       VARCHAR(255) NOT NULL,
    description VARCHAR(5000),
    status      VARCHAR(20) NOT NULL DEFAULT 'pending',
    priority    VARCHAR(10) NOT NULL DEFAULT 'medium',
    due_date    TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_todos_user_id ON todos(user_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS tags (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL,
    name       VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, name)
);

CREATE TABLE IF NOT EXISTS todo_tags (
    todo_id UUID NOT NULL REFERENCES todos(id) ON DELETE CASCADE,
    tag_id  UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (todo_id, tag_id)
);

CREATE TABLE IF NOT EXISTS search_outbox (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    todo_id    UUID NOT NULL,
    operation  VARCHAR(10) NOT NULL,
    payload    JSONB,
    processed  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_search_outbox_unprocessed ON search_outbox(processed, created_at) WHERE processed = FALSE;
