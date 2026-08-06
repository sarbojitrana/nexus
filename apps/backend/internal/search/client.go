package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/rs/zerolog"

	"github.com/sarbojitrana/nexus/internal/config"
)

const (
	IndexPosts       = "nexus_posts"
	IndexUsers       = "nexus_users"
	IndexCommunities = "nexus_communities"
)

type Client struct {
	api     *opensearchapi.Client
	enabled bool
}

func NewClient(ctx context.Context, cfg *config.Config, logger *zerolog.Logger) *Client {
	if !cfg.Search.IsConfigured() {
		return &Client{enabled: false}
	}

	api, err := opensearchapi.NewClient(opensearchapi.Config{
		Client: opensearch.Config{
			Addresses: []string{cfg.Search.URL},
			Username:  cfg.Search.Username,
			Password:  cfg.Search.Password,
		},
	})
	if err != nil {
		logger.Warn().Err(err).Msg("failed to create opensearch client, search disabled")
		return &Client{enabled: false}
	}

	c := &Client{api: api, enabled: true}

	if err := c.ensureIndices(ctx); err != nil {
		logger.Warn().Err(err).Msg("failed to ensure opensearch indices, search disabled")
		return &Client{enabled: false}
	}

	return c
}

func (c *Client) Enabled() bool {
	return c != nil && c.enabled
}

var indexMappings = map[string]string{
	IndexPosts: `{"mappings":{"properties":{
		"authorId":{"type":"keyword"},
		"communityId":{"type":"keyword"},
		"title":{"type":"text"},
		"content":{"type":"text"},
		"upvotes":{"type":"integer"},
		"commentCount":{"type":"integer"},
		"createdAt":{"type":"date"}
	}}}`,
	IndexUsers: `{"mappings":{"properties":{
		"username":{"type":"text"},
		"displayName":{"type":"text"},
		"bio":{"type":"text"},
		"followerCount":{"type":"integer"},
		"createdAt":{"type":"date"}
	}}}`,
	IndexCommunities: `{"mappings":{"properties":{
		"name":{"type":"text"},
		"slug":{"type":"keyword"},
		"description":{"type":"text"},
		"membersCount":{"type":"integer"},
		"createdAt":{"type":"date"}
	}}}`,
}

func (c *Client) ensureIndices(ctx context.Context) error {
	for index, mapping := range indexMappings {
		_, err := c.api.Indices.Create(ctx, opensearchapi.IndicesCreateReq{
			Index: index,
			Body:  strings.NewReader(mapping),
		})
		if err != nil && !strings.Contains(err.Error(), "resource_already_exists_exception") {
			return fmt.Errorf("failed to create index %s: %w", index, err)
		}
	}
	return nil
}

func (c *Client) index(ctx context.Context, index, id string, doc any) error {
	if !c.Enabled() {
		return nil
	}

	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	_, err = c.api.Index(ctx, opensearchapi.IndexReq{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(body),
	})
	if err != nil {
		return fmt.Errorf("failed to index document %s/%s: %w", index, id, err)
	}
	return nil
}

func (c *Client) delete(ctx context.Context, index, id string) error {
	if !c.Enabled() {
		return nil
	}

	_, err := c.api.Document.Delete(ctx, opensearchapi.DocumentDeleteReq{Index: index, DocumentID: id})
	if err != nil && !strings.Contains(err.Error(), "not_found") {
		return fmt.Errorf("failed to delete document %s/%s: %w", index, id, err)
	}
	return nil
}

type PostDoc struct {
	ID           uuid.UUID  `json:"id"`
	AuthorID     string     `json:"authorId"`
	CommunityID  *uuid.UUID `json:"communityId"`
	Title        *string    `json:"title"`
	Content      *string    `json:"content"`
	Upvotes      int        `json:"upvotes"`
	CommentCount int        `json:"commentCount"`
	CreatedAt    time.Time  `json:"createdAt"`
}

func (c *Client) IndexPost(ctx context.Context, doc PostDoc) error {
	return c.index(ctx, IndexPosts, doc.ID.String(), doc)
}

func (c *Client) DeletePost(ctx context.Context, id uuid.UUID) error {
	return c.delete(ctx, IndexPosts, id.String())
}

type UserDoc struct {
	ID            string    `json:"id"`
	Username      string    `json:"username"`
	DisplayName   string    `json:"displayName"`
	Bio           *string   `json:"bio"`
	FollowerCount int       `json:"followerCount"`
	CreatedAt     time.Time `json:"createdAt"`
}

func (c *Client) IndexUser(ctx context.Context, doc UserDoc) error {
	return c.index(ctx, IndexUsers, doc.ID, doc)
}

func (c *Client) DeleteUser(ctx context.Context, id string) error {
	return c.delete(ctx, IndexUsers, id)
}

type CommunityDoc struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  *string   `json:"description"`
	MembersCount int       `json:"membersCount"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (c *Client) IndexCommunity(ctx context.Context, doc CommunityDoc) error {
	return c.index(ctx, IndexCommunities, doc.ID.String(), doc)
}

func (c *Client) DeleteCommunity(ctx context.Context, id uuid.UUID) error {
	return c.delete(ctx, IndexCommunities, id.String())
}

type Results struct {
	Posts       []PostDoc      `json:"posts"`
	Users       []UserDoc      `json:"users"`
	Communities []CommunityDoc `json:"communities"`
}

const resultsPerType = 5

func (c *Client) Search(ctx context.Context, query string) (*Results, error) {
	results := &Results{Posts: []PostDoc{}, Users: []UserDoc{}, Communities: []CommunityDoc{}}
	if !c.Enabled() || strings.TrimSpace(query) == "" {
		return results, nil
	}

	var wg sync.WaitGroup
	var firstErr error
	var mu sync.Mutex

	search := func(index string, fields []string, into func(json.RawMessage) error) {
		defer wg.Done()

		body, err := json.Marshal(map[string]any{
			"query": map[string]any{
				"multi_match": map[string]any{
					"query":  query,
					"fields": fields,
				},
			},
			"size": resultsPerType,
		})
		if err != nil {
			return
		}

		resp, err := c.api.Search(ctx, &opensearchapi.SearchReq{
			Indices: []string{index},
			Body:    bytes.NewReader(body),
		})
		if err != nil {
			mu.Lock()
			if firstErr == nil {
				firstErr = fmt.Errorf("search %s: %w", index, err)
			}
			mu.Unlock()
			return
		}

		for _, hit := range resp.Hits.Hits {
			if err := into(hit.Source); err != nil {
				continue
			}
		}
	}

	wg.Add(3)
	go search(IndexPosts, []string{"title", "content"}, func(src json.RawMessage) error {
		var doc PostDoc
		if err := json.Unmarshal(src, &doc); err != nil {
			return err
		}
		mu.Lock()
		results.Posts = append(results.Posts, doc)
		mu.Unlock()
		return nil
	})
	go search(IndexUsers, []string{"username", "displayName", "bio"}, func(src json.RawMessage) error {
		var doc UserDoc
		if err := json.Unmarshal(src, &doc); err != nil {
			return err
		}
		mu.Lock()
		results.Users = append(results.Users, doc)
		mu.Unlock()
		return nil
	})
	go search(IndexCommunities, []string{"name", "description"}, func(src json.RawMessage) error {
		var doc CommunityDoc
		if err := json.Unmarshal(src, &doc); err != nil {
			return err
		}
		mu.Lock()
		results.Communities = append(results.Communities, doc)
		mu.Unlock()
		return nil
	})
	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}
	return results, nil
}
