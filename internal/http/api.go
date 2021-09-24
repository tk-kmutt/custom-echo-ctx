package api

import (
	"custom-echo-ctx/internal/http/gen"
	"custom-echo-ctx/internal/http/usecase"
	"custom-echo-ctx/pkg/context"

	"github.com/labstack/echo/v4"
)

type Api struct {
	user *usecase.User
}

func wrap(h func(c *context.Context) error, c echo.Context) error {
	return h(c.(*context.Context))
}

func NewApi() *Api {
	return &Api{user: usecase.NewUser()}
}

var _ gen.ServerInterface = (*Api)(nil)

func (p Api) Login(ctx echo.Context) error {
	return wrap(p.user.Login, ctx)
}

func (p Api) Signup(ctx echo.Context) error {
	return wrap(p.user.Signup, ctx)
}
