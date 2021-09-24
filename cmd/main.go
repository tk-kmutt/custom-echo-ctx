package main

import (
	"custom-echo-ctx/infra/mysql/repository"
	api "custom-echo-ctx/internal/http"
	"custom-echo-ctx/internal/http/gen"
	"fmt"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	om "github.com/deepmap/oapi-codegen/pkg/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"
)

type AlpacaUser struct {
	Name  string `json:"name" form:"name" query:"name" validate:"required"`
	Email string `json:"email" form:"email" query:"email" validate:"required"`
}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

// Context
// echo.Context をラップする構造体を定義する
type Context struct {
	echo.Context
	user *AlpacaUser
}

// BindValidate
// Bind と Validate を合わせたメソッド
func (c *Context) BindValidate(i interface{}) error {
	if err := c.Bind(i); err != nil {
		return c.String(http.StatusBadRequest, "Request is failed: "+err.Error())
	}
	if err := c.Validate(i); err != nil {
		return c.String(http.StatusBadRequest, "Validate is failed: "+err.Error())
	}
	return nil
}

// LogBindValidate
// Log とBind と Validate を合わせたメソッド
// funcを生やす練習
func (c *Context) LogBindValidate(i interface{}) error {
	c.Logger().Print(c.user)

	if err := c.Bind(i); err != nil {
		return c.String(http.StatusBadRequest, "Request is failed: "+err.Error())
	}
	if err := c.Validate(i); err != nil {
		return c.String(http.StatusBadRequest, "Validate is failed: "+err.Error())
	}
	return nil
}

type callFunc func(c *Context) error

func c(h callFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return h(c.(*Context))
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Validator = &Validator{validator: validator.New()}

	// validator
	spec, err := gen.GetSwagger()
	if err != nil {
		panic(err)
	}
	e.Use(om.OapiRequestValidator(spec))

	// echo.Context をラップして扱うために middleware として登録する
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return h(&Context{
				Context: c,
				user: &AlpacaUser{
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

	e.POST("/post_profile", c(func(c *Context) error {
		u := new(AlpacaUser)
		if err := c.LogBindValidate(u); err != nil {
			return err
		}
		fmt.Println(u)
		return c.String(http.StatusOK, "OK")
	}))
	gen.RegisterHandlers(e, api.NewApi(db))

	e.Logger.Fatal(e.Start(":3000"))
}
