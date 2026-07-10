CREATE TABLE
    community_reports (
        id UUID PRIMARY KEY gen_random_uuid (),
        post_id UUID NOT NULl REFERENCES posts ON DELETE CASCADE,
        reporter_id UUID NOT NULl REFERENCES users ON DELETE CASCADE,
        community_id NOT NULL REFERENCES communities ON DELETE CASCADE,
        reason TEXT NOT NULl,
        status TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX community_reports_post_id ON community_reports(post_id);
CREATE INDEX community_reports_reporter_id ON community_reports(reporter_id);
CREATE INDEX community_reports_community_id ON community_reports(community_id);
