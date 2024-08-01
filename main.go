package main

import (
	"github.com/labstack/echo/v4"
	"github.com/anakilang-ai/backend/routes"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		routes.URL(c.Response().Writer, c.Request())
		return nil
	})
	port := ":8080"
	e.Logger.Fatal(e.Start(port))
}
