ALTER TABLE users
ADD COLUMN profile_visibility TEXT NOT NULL DEFAULT 'public' CHECK (profile_visibility IN ('public', 'followers_only', 'private')),
ADD COLUMN show_online_status BOOLEAN NOT NULL DEFAULT true,
ADD COLUMN group_invite_permission TEXT NOT NULL DEFAULT 'everyone' CHECK (group_invite_permission IN ('everyone', 'followers_only', 'no_one')),
ADD COLUMN share_read_receipts BOOLEAN NOT NULL DEFAULT true;
