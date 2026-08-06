CREATE TABLE
    conversations (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        is_group BOOLEAN NOT NULL DEFAULT false,
        name TEXT,
        created_by TEXT NOT NULL REFERENCES users ON DELETE CASCADE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE TRIGGER set_updated_at_conversations BEFORE
UPDATE ON conversations FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at ();

CREATE TABLE
    conversation_participants (
        conversation_id UUID NOT NULL REFERENCES conversations ON DELETE CASCADE,
        user_id TEXT NOT NULL REFERENCES users ON DELETE CASCADE,
        joined_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        last_read_at TIMESTAMPTZ,
        PRIMARY KEY (conversation_id, user_id)
    );

CREATE INDEX idx_conversation_participants_user_id ON conversation_participants (user_id);

CREATE TABLE
    conversation_invitations (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        conversation_id UUID NOT NULL REFERENCES conversations ON DELETE CASCADE,
        inviter_id TEXT NOT NULL REFERENCES users ON DELETE CASCADE,
        invitee_id TEXT NOT NULL REFERENCES users ON DELETE CASCADE,
        status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'rejected')),
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        responded_at TIMESTAMPTZ,
        UNIQUE (conversation_id, invitee_id)
    );

CREATE INDEX idx_conversation_invitations_invitee_id ON conversation_invitations (invitee_id, status);

CREATE TABLE
    messages (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        conversation_id UUID NOT NULL REFERENCES conversations ON DELETE CASCADE,
        sender_id TEXT NOT NULL REFERENCES users ON DELETE CASCADE,
        content TEXT,
        attachment_key TEXT,
        attachment_mime_type TEXT,
        attachment_file_size BIGINT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT message_has_content_or_attachment CHECK (
            content IS NOT NULL
            OR attachment_key IS NOT NULL
        )
    );

CREATE INDEX idx_messages_conversation_id_created_at ON messages (conversation_id, created_at DESC);

CREATE OR REPLACE FUNCTION touch_conversation_on_message () RETURNS TRIGGER AS $$
BEGIN
	UPDATE conversations SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.conversation_id;
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_messages_touch_conversation
AFTER INSERT ON messages FOR EACH ROW
EXECUTE FUNCTION touch_conversation_on_message ();
