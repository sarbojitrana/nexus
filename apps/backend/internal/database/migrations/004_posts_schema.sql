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
        comment_count INT NOT NULL DEFAULT 0,
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
        vote_type TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_post_votes_post_id ON post_votes (post_id);

CREATE TRIGGER set_updated_at_post_votes BEFORE
UPDATE ON post_votes FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at ();

CREATE OR REPLACE FUNCTION increment_ancestor_comment_counts() RETURNS TRIGGER AS $$
BEGIN
	IF NEW.parent_post_id IS NOT NULL THEN
		WITH RECURSIVE ancestors AS (
			SELECT id, parent_post_id FROM posts WHERE id = NEW.parent_post_id
			UNION ALL
			SELECT p.id, p.parent_post_id
			FROM posts p
			JOIN ancestors a ON p.id = a.parent_post_id
		)
		UPDATE posts
		SET comment_count = comment_count + 1
		WHERE id IN (SELECT id FROM ancestors);
	END IF;
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_posts_ancestor_comment_count
AFTER INSERT ON posts
FOR EACH ROW EXECUTE FUNCTION increment_ancestor_comment_counts();



CREATE OR REPLACE FUNCTION sync_posts_count() RETURNS TRIGGER AS $$
BEGIN
	IF TG_OP = 'INSERT' THEN
		IF NEW.post_type = 'post' THEN
			UPDATE users SET posts_count = posts_count + 1 WHERE id = NEW.author_id;
			IF NEW.community_id IS NOT NULL THEN
				UPDATE communities SET posts_count = posts_count + 1 WHERE id = NEW.community_id;
			END IF;
		END IF;
	ELSIF TG_OP = 'DELETE' THEN
		IF OLD.post_type = 'post' THEN
			UPDATE users SET posts_count = posts_count - 1 WHERE id = OLD.author_id;
			IF OLD.community_id IS NOT NULL THEN
				UPDATE communities SET posts_count = posts_count - 1 WHERE id = OLD.community_id;
			END IF;
		END IF;
	END IF;
	RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_posts_count
AFTER INSERT OR DELETE ON posts
FOR EACH ROW EXECUTE FUNCTION sync_posts_count();