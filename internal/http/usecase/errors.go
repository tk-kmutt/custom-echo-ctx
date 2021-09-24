package usecase

import (
	"custom-echo-ctx/internal/http/gen"

	"github.com/labstack/echo/v4"
)

func sendError(ctx echo.Context, code int, message string) error {
	e := gen.Error{
		Code:    int32(code),
		Message: message,
	}
	err := ctx.JSON(code, e)
	return err
}
