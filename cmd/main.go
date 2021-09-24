package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"
)

type User struct {
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
	user *User
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
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Validator = &Validator{validator: validator.New()}

	// echo.Context をラップして扱うために middleware として登録する
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			return h(&Context{
				Context: c,
				user: &User{
					Name:  "john",
					Email: "john@cayto.jp",
				},
			})
		}
	})

	e.POST("/post_profile", c(func(c *Context) error {
		u := new(User)
		if err := c.LogBindValidate(u); err != nil {
			return err
		}
		fmt.Println(u)
		return c.String(http.StatusOK, "OK")
	}))

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}
