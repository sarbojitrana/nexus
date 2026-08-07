-- The Post model and the trending-feed ranking query both reference
-- posts.engagement and posts.is_popularity_updated, but neither column was
-- ever created. Every post read path (GET /feed, GET /posts/:id) failed with
-- a scan/column error until these existed.
--
-- engagement is derived, not written by application code, so it's a stored
-- generated column: Postgres keeps it in sync with the vote/comment counters
-- that existing triggers already maintain.
ALTER TABLE posts
    ADD COLUMN engagement INT GENERATED ALWAYS AS (upvotes + comment_count) STORED;

ALTER TABLE posts
    ADD COLUMN is_popularity_updated BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX idx_posts_engagement_created_at ON posts (engagement DESC, created_at DESC);
