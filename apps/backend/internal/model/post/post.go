package post

import (
	"time"

	"github.com/google/uuid"
	"github.com/sarbojitrana/nexus/internal/model"
)

type PostType string

const (
	PostTypeComment PostType = "comment"
	PostTypePost    PostType = "post"
)

type VoteType int16

const (
	Upvote   VoteType = 1
	DownVote VoteType = -1
)

type Post struct {
	model.Base
	AuthorID     uuid.UUID  `json:"authorId" db:"author_id"`
	CommunityID  *uuid.UUID `json:"communityId" db:"community_id"`
	ParentPostID *uuid.UUID `json:"parentPostId" db:"parent_post_id"`
	PostType     PostType   `json:"postType" db:"post_type"`
	Title        *string    `json:"title" db:"title"`
	Content      *string    `json:"content" db:"content"`
	Upvotes      int        `json:"upvotes" db:"upvotes"`
	Downvotes    int        `json:"downvotes" db:"downvotes"`
	CommentCount int        `json:"commentCount" db:"comment_count"`
	DeletedAt    *time.Time `json:"deletedAt" db:"deleted_at"`
}

type PostMedia struct {
	model.BaseWithId
	UserID      uuid.UUID `json:"userId" db:"user_id"`
	PostID      uuid.UUID `json:"postId" db:"post_id"`
	DownloadKey string    `json:"downloadKey" db:"download_key"`
	FileSize    int64     `json:"fileSize" db:"file_size"`
	MimeType    string    `json:"mimeType" db:"mime_type"`
	model.BaseWithCreatedAt
}

type PostVote struct {
	PostID   uuid.UUID `json:"postId" db:"post_id"`
	UserID   uuid.UUID `json:"userId" db:"user_id"`
	VoteType VoteType  `json:"voteType" db:"vote_type"`
	model.BaseWithCreatedAt
	model.BaseWithUpdatedAt
}

type PopulatedPost struct {
	Post
	PostMedia []PostMedia `json:"postMedia" db:"post_media"`
}
