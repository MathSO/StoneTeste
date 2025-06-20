package main

import (
	"server/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	tickerHandlers := handlers.Ticker{}
	e.GET("/info/:ticker", tickerHandlers.GetInfo)

	e.Start(":8080")
}
