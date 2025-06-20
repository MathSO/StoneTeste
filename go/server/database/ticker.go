package database

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
)

type Ticker struct{}

func (Ticker) GetInfo(ticker string, date string) (info TickerGetInfo, err error) {
	conn, err := connect()
	if err != nil {
		return
	}

	var row pgx.Row
	if date != "" {
		sql := `
		WITH aux as (
			SELECT 
				SUM(preco_negocio) AS sum_preco_negocio, 
				MAX(preco_negocio) AS max_preco_negocio
			FROM negociacoes n 
			WHERE 
				codigo_instrumento = $1 
				AND data_negocio >= DATE($2)
			GROUP BY data_negocio
		) SELECT max(sum_preco_negocio), max(max_preco_negocio) FROM aux;`
		row = conn.QueryRow(context.Background(), sql, ticker, date)
	} else {
		sql := `
		WITH aux as (
			SELECT 
				SUM(preco_negocio) AS sum_preco_negocio, 
				MAX(preco_negocio) AS max_preco_negocio
			FROM negociacoes n 
			WHERE 
				codigo_instrumento = $1
			GROUP BY data_negocio
		) SELECT max(sum_preco_negocio), max(max_preco_negocio) FROM aux;`
		row = conn.QueryRow(context.Background(), sql, ticker)
	}

	err = row.Scan(&info.MaxDailyVolume, &info.MaxRangeValue)
	info.Ticker = ticker
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
	}

	return
}
