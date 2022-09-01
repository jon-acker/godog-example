package main

import (
	"context"
	echo "github.com/labstack/echo/v4" //nolint:unused
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

type Application struct {
	Router   http.Handler
	Database Database
}
type Database struct {
	Library Library
	Loans   map[string]string
}

type Library struct {
	Members []string
}

func (l *Library) HasMember(memberName string) bool {
	for _, member := range l.Members {
		if memberName == member {
			return true
		}
	}

	return false
}

func NewApplication() *Application {
	app := &Application{
		Database: Database{
			Library: Library{
				Members: []string{},
			},
			Loans: make(map[string]string),
		},
	}

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:2017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer cancel()
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

type BorrowRequest struct {
	MemberName string `json:"member_name"`
	BookName   string `json:"book_name"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (a *Application) register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// member gets registered
	a.Database.Library.Members = append(a.Database.Library.Members, req.MemberName)
	return c.JSON(http.StatusCreated, &req)
}

func (a *Application) borrow(c echo.Context) error {
	var req BorrowRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if a.Database.Library.HasMember(req.MemberName) {
		a.Database.Loans[req.MemberName] = req.BookName
		return c.JSON(http.StatusCreated, "")
	}

	var response ErrorResponse
	response.Message = "Only members can borrow books"
	return c.JSON(http.StatusUnauthorized, &response)
}
