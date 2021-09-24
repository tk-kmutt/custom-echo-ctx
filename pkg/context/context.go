package context

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

// Context
// echo.Context をラップする構造体を定義する
type Context struct {
	echo.Context
	User *Auth
	DB   *gorm.DB
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
	c.Logger().Print(c.User)

	if err := c.Bind(i); err != nil {
		return c.String(http.StatusBadRequest, "Request is failed: "+err.Error())
	}
	if err := c.Validate(i); err != nil {
		return c.String(http.StatusBadRequest, "Validate is failed: "+err.Error())
	}
	return nil
}

type callFunc func(c *Context) error

func Wrap(h callFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return h(c.(*Context))
	}
}
