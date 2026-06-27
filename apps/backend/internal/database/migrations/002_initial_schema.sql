CREATE TABLE
    users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
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

CREATE UNIQUE INDEX idx_users_username ON users (username);

CREATE TRIGGER set_updated_at_users BEFORE
UPDATE ON users FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at ();

CREATE TABLE
    user_follows (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        follower_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
        following_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT self_follow_check CHECK (follower_id <> following_id)
    );

CREATE TRIGGER set_updated_at_user_follows BEFORE
UPDATE ON user_follows FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at ();

CREATE TABLE
    communities (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        admin_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
        name VARCHAR(50) NOT NULL,
        slug TEXT UNIQUE NOT NULL,
        description TEXT,
        avatar_key TEXT,
        banner_key TEXT,
        members_count INT NOT NULL DEFAULT 0,
        posts_count INT NOT NULL DEFAULT 0,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE UNIQUE INDEX idx_communities_name ON communities (name);

CREATE INDEX idx_communities_owner_id ON communities (admin_id);

CREATE TRIGGER set_updated_at_communities BEFORE
UPDATE ON communities FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at ();

CREATE TABLE
    community_follows (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        follower_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
        community_id UUID NOT NULL REFERENCES communities ON DELETE CASCADE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE TRIGGER set_updated_at_community_follows BEFORE
UPDATE ON communtiy_follows FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at ();

CREATE TABLE
    community_members (
        user_id UUID REFERENCES users ON DELETE CASCADE,
        community_id UUID REFERENCES communities ON DELETE CASCADE,
        PRIMARY KEY (user_id, community_id),
        role TEXT NOT NULL DEFAULT 'member',
        joined_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_community_members_user_id ON community_members (user_id);

CREATE TRIGGER set_updated_at_community_members BEFORE
UPDATE ON community_members FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at ();

CREATE TABLE
    posts (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        author_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
        community_id UUID REFERENCES communities ON DELETE CASCADE,
        parent_post_id UUID REFERENCES posts ON DELETE CASCADE,
        post_type VARCHAR(50) NOT NULL DEFAULT 'post',
        title VARCHAR(100),
        content TEXT,
        upvotes INT NOT NULL DEFAULT 0,
        downvotes INT NOT NULL DEFAULT 0,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        deleted_at TIMESTAMPTZ
    );

CREATE INDEX idx_posts_author_id ON posts (author_id);

CREATE TRIGGER set_updated_at_posts BEFORE
UPDATE ON posts FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at ();

CREATE TABLE
    post_media (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        post_id UUID NOT NULL REFERENCES posts ON DELETE CASCADE,
        download_key TEXT NOT NULL,
        file_size BIGINT NOT NULL,
        mime_type TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_post_media_post_id ON post_media (post_id);

CREATE TABLE
    post_votes (
        post_id UUID NOT NULL REFERENCES posts ON DELETE CASCADE,
        user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
        PRIMARY KEY (post_id, user_id),
        vote_type SMALLINT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_post_votes_post_id ON post_votes (post_id);

CREATE TRIGGER set_updated_at_post_votes BEFORE
UPDATE ON post_votes FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at ();