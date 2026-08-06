CREATE TABLE
    notifications (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        user_id TEXT NOT NULL REFERENCES users ON DELETE CASCADE,
        actor_id TEXT REFERENCES users ON DELETE SET NULL,
        type TEXT NOT NULL CHECK (type IN ('follow', 'group_invitation', 'message')),
        data JSONB NOT NULL DEFAULT '{}',
        is_read BOOLEAN NOT NULL DEFAULT false,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_notifications_user_id_created_at ON notifications (user_id, created_at DESC);
