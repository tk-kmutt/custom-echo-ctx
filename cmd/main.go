package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	validator "gopkg.in/go-playground/validator.v9"
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

func main() {
	e := echo.New()
	e.Validator = &Validator{validator: validator.New()}

	e.POST("/post_profile", func(c echo.Context) error {
		u := new(User)
		if err := c.Bind(u); err != nil {
			return c.String(http.StatusBadRequest, "Request is failed: "+err.Error())
		}
		if err := c.Validate(u); err != nil {
			return c.String(http.StatusBadRequest, "Validate is failed: "+err.Error())
		}
		fmt.Println(u)
		return c.String(http.StatusOK, "OK")
	})

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}
