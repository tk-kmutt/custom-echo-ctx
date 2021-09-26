package internal

import (
	"context"
	"custom-echo-ctx/infra/mysql/repository"
	api "custom-echo-ctx/internal/http"
	"custom-echo-ctx/internal/http/gen"
	"custom-echo-ctx/pkg/jwt"
	"errors"
	"fmt"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/getkin/kin-openapi/openapi3filter"

	om "github.com/deepmap/oapi-codegen/pkg/middleware"

	mc "custom-echo-ctx/pkg/context"
	mv "custom-echo-ctx/pkg/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

type AlpacaUser struct {
	Name  string `json:"name" form:"name" query:"name" validate:"required"`
	Email string `json:"email" form:"email" query:"email" validate:"required"`
}

func (r *Server) Run() {
	e := echo.New()
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Validator = &mv.Validator{Validator: validator.New()}

	// mysql connection
	dsn := "user:pass@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	if err := db.AutoMigrate(&repository.User{}); err != nil {
		panic(err.Error())
	}

	// echo.Context をラップして扱うために middleware として登録する
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return h(&mc.Context{
				Context: c,
				Auth: &jwt.LoginUser{
					Name:  "john",
					Email: "john@cayto.jp",
				},
				DB: db,
			})
		}
	})

	e.POST("/post_profile", mc.Wrap(func(c *mc.Context) error {
		u := new(AlpacaUser)
		if err := c.LogBindValidate(u); err != nil {
			return err
		}
		fmt.Println(u)
		return c.String(http.StatusOK, "OK")
	}))

	// validator
	spec, err := gen.GetSwagger()
	if err != nil {
		panic(err)
	}
	validatorOptions := &om.Options{}
	validatorOptions.Options.AuthenticationFunc = func(c context.Context, input *openapi3filter.AuthenticationInput) error {
		h := input.RequestValidationInput.Request.Header["Authorization"]
		if h == nil {
			return errors.New("auth failed")
		}

		if h[0] != "Bearer super_strong_password" {
			return errors.New("auth failed")
		}
		return nil
	}
	oapi := e.Group("")
	oapi.Use(om.OapiRequestValidatorWithOptions(spec, validatorOptions))
	gen.RegisterHandlers(oapi, api.NewApi())

	e.Logger.Fatal(e.Start(":3000"))
}
