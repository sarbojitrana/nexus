package repository

import "github.com/sarbojitrana/nexus/internal/server"

type Repositories struct{}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{}
}
