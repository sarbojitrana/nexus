package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/sarbojitrana/nexus/internal/search"
	"github.com/sarbojitrana/nexus/internal/server"
)

// SearchRepository is the Postgres fallback behind SearchService. It trades
// relevance ranking and typo tolerance for always being available: OpenSearch
// can be unconfigured, unreachable, or holding a stale/empty index, and search
// still has to return something useful from the source of truth.
type SearchRepository struct {
	server *server.Server
}

func NewSearchRepository(server *server.Server) *SearchRepository {
	return &SearchRepository{server: server}
}

const fallbackLimit = 5

func (r *SearchRepository) SearchPosts(ctx context.Context, query string) ([]search.PostDoc, error) {
	rows, err := r.server.DB.Pool.Query(ctx, `
		SELECT id, author_id, community_id, title, content, upvotes, comment_count, created_at
		FROM posts
		WHERE deleted_at IS NULL
			AND post_type = 'post'
			AND (title ILIKE @pattern OR content ILIKE @pattern)
		ORDER BY (upvotes + comment_count) DESC, created_at DESC
		LIMIT @limit
	`, pgx.NamedArgs{"pattern": "%" + query + "%", "limit": fallbackLimit})
	if err != nil {
		return nil, fmt.Errorf("failed to search posts: %w", err)
	}
	defer rows.Close()

	docs := []search.PostDoc{}
	for rows.Next() {
		var d search.PostDoc
		if err := rows.Scan(&d.ID, &d.AuthorID, &d.CommunityID, &d.Title, &d.Content,
			&d.Upvotes, &d.CommentCount, &d.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan post search row: %w", err)
		}
		docs = append(docs, d)
	}
	return docs, rows.Err()
}

func (r *SearchRepository) SearchUsers(ctx context.Context, query string) ([]search.UserDoc, error) {
	rows, err := r.server.DB.Pool.Query(ctx, `
		SELECT id, username, display_name, bio, follower_count, created_at
		FROM users
		WHERE profile_visibility = 'public'
			AND (username ILIKE @pattern OR display_name ILIKE @pattern)
		ORDER BY follower_count DESC, created_at DESC
		LIMIT @limit
	`, pgx.NamedArgs{"pattern": "%" + query + "%", "limit": fallbackLimit})
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	docs := []search.UserDoc{}
	for rows.Next() {
		var d search.UserDoc
		if err := rows.Scan(&d.ID, &d.Username, &d.DisplayName, &d.Bio,
			&d.FollowerCount, &d.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user search row: %w", err)
		}
		docs = append(docs, d)
	}
	return docs, rows.Err()
}

func (r *SearchRepository) SearchCommunities(ctx context.Context, query string) ([]search.CommunityDoc, error) {
	rows, err := r.server.DB.Pool.Query(ctx, `
		SELECT id, name, slug, description, members_count, created_at
		FROM communities
		WHERE name ILIKE @pattern OR slug ILIKE @pattern OR description ILIKE @pattern
		ORDER BY members_count DESC, created_at DESC
		LIMIT @limit
	`, pgx.NamedArgs{"pattern": "%" + query + "%", "limit": fallbackLimit})
	if err != nil {
		return nil, fmt.Errorf("failed to search communities: %w", err)
	}
	defer rows.Close()

	docs := []search.CommunityDoc{}
	for rows.Next() {
		var d search.CommunityDoc
		if err := rows.Scan(&d.ID, &d.Name, &d.Slug, &d.Description,
			&d.MembersCount, &d.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan community search row: %w", err)
		}
		docs = append(docs, d)
	}
	return docs, rows.Err()
}
