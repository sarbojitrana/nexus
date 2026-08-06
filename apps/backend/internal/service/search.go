package service

import (
	"context"

	"github.com/sarbojitrana/nexus/internal/middleware"
	"github.com/sarbojitrana/nexus/internal/search"
	"github.com/sarbojitrana/nexus/internal/server"
)

type SearchService struct {
	server *server.Server
	search *search.Client
}

func NewSearchService(s *server.Server, search *search.Client) *SearchService {
	return &SearchService{server: s, search: search}
}

func (s *SearchService) Search(ctx context.Context, query string) (*search.Results, error) {
	results, err := s.search.Search(ctx, query)
	if err != nil {
		middleware.GetLoggerFromContext(ctx).Error().Err(err).Str("query", query).Msg("search failed")
		return nil, err
	}
	return results, nil
}
