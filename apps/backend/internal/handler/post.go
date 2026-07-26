package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/sarbojitrana/nexus/internal/middleware"
	"github.com/sarbojitrana/nexus/internal/model"
	"github.com/sarbojitrana/nexus/internal/model/post"
	"github.com/sarbojitrana/nexus/internal/server"
	"github.com/sarbojitrana/nexus/internal/service"
)

type PostHandler struct {
	Handler
	services *service.Services
}

func NewPostHandler(s *server.Server, services *service.Services) *PostHandler {
	return &PostHandler{
		Handler:  NewHandler(s),
		services: services,
	}
}

func (h *PostHandler) CreatePost(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *post.CreatePostPayload) (*post.Post, error) {
		return h.services.Post.CreatePost(c.Request().Context(), middleware.GetUserID(c), req)
	}, http.StatusCreated, &post.CreatePostPayload{})(c)
}

func (h *PostHandler) UpdatePost(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *post.UpdatePostByIDPayload) (*post.Post, error) {
		postID, err := parseUUIDParam(c, "id")
		if err != nil {
			return nil, err
		}
		return h.services.Post.UpdatePost(c.Request().Context(), middleware.GetUserID(c), postID, req)
	}, http.StatusOK, &post.UpdatePostByIDPayload{})(c)
}

func (h *PostHandler) DeletePost(c echo.Context) error {
	return HandleNoContent(h.Handler, func(c echo.Context, req *model.Empty) error {
		postID, err := parseUUIDParam(c, "id")
		if err != nil {
			return err
		}
		return h.services.Post.DeletePost(c.Request().Context(), middleware.GetUserID(c), postID)
	}, http.StatusNoContent, &model.Empty{})(c)
}

func (h *PostHandler) GetPostByID(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *model.Empty) (*post.PopulatedPost, error) {
		postID, err := parseUUIDParam(c, "id")
		if err != nil {
			return nil, err
		}
		return h.services.Post.GetPostByID(c.Request().Context(), viewerIDFromContext(c), postID)
	}, http.StatusOK, &model.Empty{})(c)
}

func (h *PostHandler) GetComments(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *post.GetCommentsByPostIDQuery) (*model.CursorPaginatedResponse[post.PopulatedPost], error) {
		postID, err := parseUUIDParam(c, "id")
		if err != nil {
			return nil, err
		}
		return h.services.Post.GetComments(c.Request().Context(), viewerIDFromContext(c), postID, req)
	}, http.StatusOK, &post.GetCommentsByPostIDQuery{})(c)
}

func (h *PostHandler) GetReplies(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *post.GetRepliesByCommentIDQuery) (*model.OffsetPaginatedResponse[post.PopulatedPost], error) {
		commentID, err := parseUUIDParam(c, "id")
		if err != nil {
			return nil, err
		}
		return h.services.Post.GetReplies(c.Request().Context(), viewerIDFromContext(c), commentID, req)
	}, http.StatusOK, &post.GetRepliesByCommentIDQuery{})(c)
}

func (h *PostHandler) React(c echo.Context) error {
	return HandleNoContent(h.Handler, func(c echo.Context, req *post.ReactToPostPayload) error {
		postID, err := parseUUIDParam(c, "id")
		if err != nil {
			return err
		}
		return h.services.Post.React(c.Request().Context(), middleware.GetUserID(c), postID, req)
	}, http.StatusNoContent, &post.ReactToPostPayload{})(c)
}

func (h *PostHandler) GetFeed(c echo.Context) error {
	return Handle(h.Handler, func(c echo.Context, req *post.GetPostsQuery) (*post.GetPostsQueryResponse, error) {
		return h.services.Post.GetFeed(c.Request().Context(), viewerIDFromContext(c), req)
	}, http.StatusOK, &post.GetPostsQuery{})(c)
}
