-- A user reporting the same post twice creates indistinguishable rows in the
-- moderation queue. Collapse existing open duplicates, then stop new ones at
-- the database.
--
-- The tie-break on id matters: rows inserted in the same statement share a
-- created_at, so comparing timestamps alone leaves duplicates behind and the
-- unique index below then fails to build.
DELETE FROM community_reports a
USING community_reports b
WHERE a.reporter_id = b.reporter_id
  AND a.post_id = b.post_id
  AND a.status = 'pending'
  AND b.status = 'pending'
  AND (a.created_at > b.created_at OR (a.created_at = b.created_at AND a.id > b.id));

-- Partial, so a post can be reported again after an earlier report is resolved
-- or dismissed -- only one *open* report per person per post.
CREATE UNIQUE INDEX idx_community_reports_one_open_per_reporter_post
    ON community_reports (reporter_id, post_id)
    WHERE status = 'pending';
