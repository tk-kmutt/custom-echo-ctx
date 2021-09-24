package api

import (
	"custom-echo-ctx/internal/http/gen"
	"custom-echo-ctx/internal/http/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Api struct {
	user *usecase.User
}

func NewApi(db *gorm.DB) *Api {
	return &Api{user: usecase.NewUser(db)}
}

var _ gen.ServerInterface = (*Api)(nil)

func (p Api) Login(ctx echo.Context) error {
	return p.user.Login(ctx)
}

func (p Api) Signup(ctx echo.Context) error {
	return p.user.Signup(ctx)
}
