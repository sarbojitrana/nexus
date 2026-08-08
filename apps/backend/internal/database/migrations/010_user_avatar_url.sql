-- Clerk owns the profile picture: its Manage Account UI is the primary place
-- users change it, and it hands back a hosted URL rather than something we
-- store. Keeping that URL alongside the R2 avatar_key lets a picture set in
-- Clerk show up everywhere in Nexus instead of only on the Clerk widget.
ALTER TABLE users ADD COLUMN avatar_url TEXT;
