CREATE TABLE
    users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        email_id TEXT NOT NULL,
        clerk_id TEXT NOT NULL,
        username VARCHAR(50) NOT NULL,
        display_name VARCHAR(50) NOT NULL,
        bio TEXT,
        avatar_key TEXT,
        banner_key TEXT,
        follower_count INT NOT NULL DEFAULT 0,
        following_count INT NOT NULL DEFAULT 0,
        posts_count INT NOT NULL DEFAULT 0,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_users_clerk_id ON users (clerk_id);
CREATE INDEX idx_users_email_id ON users (email_id);
CREATE UNIQUE INDEX idx_users_username ON users (username);

CREATE TRIGGER set_updated_at_users BEFORE
UPDATE ON users FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at ();

CREATE TABLE
    user_follows (
        follower_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
        following_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
        PRIMARY KEY (follower_id, following_id),
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT self_follow_check CHECK (follower_id <> following_id)
    );
CREATE INDEX idx_user_follows_follower_id ON user_follows (follower_id);
CREATE INDEX idx_user_follows_following_id ON user_follows (following_id);



CREATE TABLE user_blocks (
    blocker_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    blocked_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (blocker_id, blocker_id)
);

CREATE INDEX user_blocks_blocker_id ON user_blocks(blocker_id);
CREATE INDEX user_blocks_blocked_id ON user_blocks(blocked_id);
