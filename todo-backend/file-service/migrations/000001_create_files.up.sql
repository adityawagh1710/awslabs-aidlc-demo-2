CREATE TABLE IF NOT EXISTS file_attachments (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    todo_id      UUID NOT NULL,
    user_id      UUID NOT NULL,
    filename     VARCHAR(255) NOT NULL,
    storage_path TEXT NOT NULL,
    mime_type    VARCHAR(100) NOT NULL,
    size_bytes   BIGINT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_file_attachments_todo ON file_attachments(todo_id);
CREATE INDEX idx_file_attachments_user ON file_attachments(user_id);
