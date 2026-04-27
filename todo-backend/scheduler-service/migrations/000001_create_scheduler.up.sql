CREATE TABLE IF NOT EXISTS reminders (
    id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    todo_id UUID NOT NULL,
    user_id UUID NOT NULL,
    fire_at TIMESTAMPTZ NOT NULL,
    fired   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_reminders_due ON reminders(fire_at) WHERE fired = FALSE;

CREATE TABLE IF NOT EXISTS recurrence_configs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    todo_id         UUID NOT NULL UNIQUE,
    cron_expression VARCHAR(100) NOT NULL,
    next_occurrence TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
