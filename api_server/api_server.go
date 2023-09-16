package api_server

import (
	"github.com/labstack/echo/v4"
	"security/api_server/handlers"
)

func StartServer() {
	e := echo.New()
	e.GET("/abc", handlers.RequestsHandler)
	e.Logger.Fatal(e.Start(":8001"))
}
