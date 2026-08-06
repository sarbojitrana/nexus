package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/sarbojitrana/nexus/internal/model/search"
	searchLib "github.com/sarbojitrana/nexus/internal/search"
	"github.com/sarbojitrana/nexus/internal/server"
	"github.com/sarbojitrana/nexus/internal/service"
)

type SearchHandler struct {
	Handler
	services *service.Services
}

func NewSearchHandler(s *server.Server, services *service.Services) *SearchHandler {
	return &SearchHandler{Handler: NewHandler(s), services: services}
}

func (h *SearchHandler) Search(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *search.SearchQuery) (*searchLib.Results, error) {
		return h.services.Search.Search(c.Request().Context(), req.Query)
	}, http.StatusOK, &search.SearchQuery{})(c)
}
