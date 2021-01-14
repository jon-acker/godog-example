package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Application struct {
	Router   http.Handler
	Database map[string]string
}

func NewApplication() *Application {
	app := &Application{
		Database: make(map[string]string),
	}
	app.newRouter()
	return app
}

func (a *Application) newRouter() {
	e := echo.New()

	e.POST("/register", a.register)
	e.PUT("/borrow", a.borrow)

	a.Router = e
}

type RegisterRequest struct {
	MemberName  string `json:"member_name"`
	LibraryName string `json:"library_name"`
}

func (a *Application) register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	// magic happens here ...
	// member gets registered
	return c.JSON(http.StatusCreated, &req)
}

func (a *Application) borrow(c echo.Context) error {
	// implement borrowing
	return c.String(http.StatusCreated, "")
}
