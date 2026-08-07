package service

import (
	"context"
	"strings"

	"github.com/sarbojitrana/nexus/internal/middleware"
	"github.com/sarbojitrana/nexus/internal/repository"
	"github.com/sarbojitrana/nexus/internal/search"
	"github.com/sarbojitrana/nexus/internal/server"
)

type SearchService struct {
	server *server.Server
	search *search.Client
	repo   *repository.SearchRepository
}

func NewSearchService(s *server.Server, searchClient *search.Client, repo *repository.SearchRepository) *SearchService {
	return &SearchService{server: s, search: searchClient, repo: repo}
}

// Search prefers OpenSearch and falls back to Postgres when it can't answer:
// unconfigured, unreachable, or holding an empty index. The empty-index case
// matters because indexing is best-effort -- a document that failed to index
// would otherwise be invisible to search forever, even though Postgres has it.
func (s *SearchService) Search(ctx context.Context, query string) (*search.Results, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	if strings.TrimSpace(query) == "" {
		return &search.Results{
			Posts:       []search.PostDoc{},
			Users:       []search.UserDoc{},
			Communities: []search.CommunityDoc{},
		}, nil
	}

	if s.search.Enabled() {
		results, err := s.search.Search(ctx, query)
		if err != nil {
			logger.Warn().Err(err).Str("query", query).Msg("opensearch failed, falling back to postgres")
		} else if len(results.Posts) > 0 || len(results.Users) > 0 || len(results.Communities) > 0 {
			return results, nil
		} else {
			logger.Debug().Str("query", query).Msg("opensearch returned nothing, checking postgres")
		}
	}

	return s.searchPostgres(ctx, query)
}

func (s *SearchService) searchPostgres(ctx context.Context, query string) (*search.Results, error) {
	logger := middleware.GetLoggerFromContext(ctx)

	results := &search.Results{
		Posts:       []search.PostDoc{},
		Users:       []search.UserDoc{},
		Communities: []search.CommunityDoc{},
	}

	// One slow lane shouldn't blank the whole response -- each is logged and
	// skipped independently.
	if posts, err := s.repo.SearchPosts(ctx, query); err != nil {
		logger.Error().Err(err).Str("query", query).Msg("postgres post search failed")
	} else {
		results.Posts = posts
	}

	if users, err := s.repo.SearchUsers(ctx, query); err != nil {
		logger.Error().Err(err).Str("query", query).Msg("postgres user search failed")
	} else {
		results.Users = users
	}

	if communities, err := s.repo.SearchCommunities(ctx, query); err != nil {
		logger.Error().Err(err).Str("query", query).Msg("postgres community search failed")
	} else {
		results.Communities = communities
	}

	return results, nil
}
