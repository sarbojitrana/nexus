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
        can_post TEXT DEFAULT 'all',
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE UNIQUE INDEX idx_communities_name ON communities (name);

CREATE INDEX idx_communities_owner_id ON communities (admin_id);

CREATE TRIGGER set_updated_at_communities BEFORE
UPDATE ON communities FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at ();

CREATE TABLE
    community_follows (
        follower_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
        community_id UUID NOT NULL REFERENCES communities ON DELETE CASCADE,
        PRIMARY KEY (follower_id, community_id),
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_community_follows_follower_id ON community_follows (follower_id);
CREATE INDEX idx_community_follows_community_id ON community_follows (community_id);

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
    communtiy_reports(
        id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        reporter_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
        community_id UUID NOT NULL REFERENCES communities ON DELETE CASCADE,
        post_id UUID NOT NULL REFERENCES posts ON DELETE CASCADE,
        reason TEXT NOT NULL,
        status VARCHAR(50) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP 
    );

CREATE INDEX idx_community_reports_reporter_id ON communtiy_reports(reporter_id);
CREATE INDEX idx_community_reports_post_id ON communtiy_reports(post_id);
CREATE INDEX idx_community_reports_community_id ON communtiy_reports(community_id);
CREATE INDEX idx_community_reports_status ON communtiy_reports(status);


CREATE TABLE banned_from_community_users (
    community_id UUID NOT NULL REFERENCES communities ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (community_id, user_id)
);

CREATE INDEX banned_from_community_users_user_id ON banned_from_community_users(user_id);
CREATE INDEX banned_from_community_users_community_id ON banned_from_community_users(community_id);

CREATE OR REPLACE FUNCTION sync_community_members_count() RETURNS TRIGGER AS $$
BEGIN
	IF TG_OP = 'INSERT' THEN
		UPDATE communities SET members_count = members_count + 1 WHERE id = NEW.community_id;
	ELSIF TG_OP = 'DELETE' THEN
		UPDATE communities SET members_count = members_count - 1 WHERE id = OLD.community_id;
	END IF;
	RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_community_members_count
AFTER INSERT OR DELETE ON community_members
FOR EACH ROW EXECUTE FUNCTION sync_community_members_count();