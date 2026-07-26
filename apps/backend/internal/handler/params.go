package handler

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/sarbojitrana/nexus/internal/errs"
)

// parseUUIDParam reads and parses a path param as a UUID, returning a clean
// 400 (rather than a raw parse error) on malformed input.
func parseUUIDParam(c echo.Context, name string) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		return uuid.Nil, errs.NewBadRequestError(fmt.Sprintf("invalid %s", name), false, nil, nil, nil)
	}
	return id, nil
}
