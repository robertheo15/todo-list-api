package main

import (
	"github.com/labstack/echo/v4/middleware"
	"todo-list-api/internal/config"
	"todo-list-api/internal/echohttp"
	"todo-list-api/internal/repository"
	"todo-list-api/internal/service"

	"github.com/labstack/echo/v4"
)

func main() {
	//time.Local = config.LocationJakarta

	config.LoadEnvFile()

	db := config.ConnectDB()
	repo := repository.NewPostgresRepository(db)
	svc := service.NewService(repo)

	e := echo.New()
	e.Debug = true
	e.Use(middleware.Logger())

	echohttp.NewServer(e, svc).Run()
}
