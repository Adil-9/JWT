package app

import (
	"fmt"
	"jwt/internal/app/database"
	"jwt/internal/app/endpoint"
	"jwt/internal/app/service"

	"github.com/gin-gonic/gin"
)

type App struct {
	db  *database.Database
	s   *service.Service
	e   *endpoint.EndPoint
	srv *gin.Engine
}

func New() (*App, error) {
	a := &App{}

	a.db = database.New()

	a.s = service.New(a.db)

	a.e = endpoint.New(a.s)

	a.srv = gin.Default()

	a.srv.GET("/access", a.e.Access)
	a.srv.GET("/refresh", a.e.Refresh)

	return a, nil
}

func (a *App) Run() error {

	fmt.Println("Server started at :8000")

	return a.srv.Run(":8000")
}
