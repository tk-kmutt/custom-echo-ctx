package main

import (
	"custom-echo-ctx/infra/mysql/repository"
	api "custom-echo-ctx/internal/http"
	"custom-echo-ctx/internal/http/gen"
	"fmt"
	"net/http"

	om "github.com/deepmap/oapi-codegen/pkg/middleware"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	mc "custom-echo-ctx/pkg/ctx"
	mv "custom-echo-ctx/pkg/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"
)

type AlpacaUser struct {
	Name  string `json:"name" form:"name" query:"name" validate:"required"`
	Email string `json:"email" form:"email" query:"email" validate:"required"`
}

func main() {
	e := echo.New()
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Validator = &mv.Validator{Validator: validator.New()}
	// echo.Context をラップして扱うために middleware として登録する
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return h(&mc.Context{
				Context: c,
				User: &mc.Auth{
					Name:  "john",
					Email: "john@cayto.jp",
				},
			})
		}
	})

	// mysql connection
	dsn := "user:pass@tcp(127.0.0.1:3306)/cec?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	if err := db.AutoMigrate(&repository.User{}); err != nil {
		panic(err.Error())
	}

	basic := e.Group("")
	oapi := e.Group("/oapi")

	basic.POST("/post_profile", mc.Wrap(func(c *mc.Context) error {
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
	oapi.Use(om.OapiRequestValidator(spec))
	gen.RegisterHandlers(oapi, api.NewApi(db))

	e.Logger.Fatal(e.Start(":3000"))
}
