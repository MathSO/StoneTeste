package handlers

import (
	"server/database"

	"github.com/labstack/echo/v4"
)

type Ticker struct{}

func (Ticker) GetInfo(ctx echo.Context) error {
	ticker := ctx.Param("ticker")
	date := ctx.QueryParam("filter_date")

	info, err := database.Ticker{}.GetInfo(ticker, date)
	if err != nil {
		return ctx.String(400, err.Error())
	}

	return ctx.JSON(200, info)
}
