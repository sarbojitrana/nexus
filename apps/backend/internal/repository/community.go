package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sarbojitrana/nexus/internal/model/post"
)




func (r *PostRepository) GetCommunityPostByID(ctx context.Context, postID uuid.UUID, communityID uuid.UUID) (*post.PopulatedPost, error) {
	stmt := `
		SELECT p.* , 
		COALESCE(
			json_agg(
				to_jsonb(camel(pm))
				ORDER BY
					pm.created_at DESC,
					pm.id DESC
			) FILTER(
				WHERE pm.id IS NOT NULL 
			),
			'[]' :: JSONB
		) AS post_media,
		json_build_object(
			'communityId', community.id
			'communityName', community.name,
			'communityAvatarKey', community.avatar_key
		) AS mini_community
		FROM posts p
		LEFT JOIN post_media pm ON pm.post_id = p.id
		LEFT JOIN communities community ON community.id = @community_id
		WHERE p.id = @post_id
		GROUP BY p.id
	`

	row, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"post_id":      postID,
		"community_id": communityID,
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to process get post by id query for post_id %s: %w", postID, err)
	}

	post, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[post.PopulatedPost])

	if err != nil {
		return nil, fmt.Errorf("Failed to parse the row to the struct: %w", err)
	}

	return &post, nil

}
