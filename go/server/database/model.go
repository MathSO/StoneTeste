package database

import "encoding/json"

type TickerGetInfo struct {
	Ticker         string
	MaxRangeValue  *float64
	MaxDailyVolume *float64
}

func (info TickerGetInfo) MarshalJSON() ([]byte, error) {
	type i TickerGetInfo
	aux := i(info)

	if aux.MaxRangeValue == nil {
		aux.MaxRangeValue = new(float64)
	}
	if aux.MaxDailyVolume == nil {
		aux.MaxDailyVolume = new(float64)
	}

	return json.Marshal(aux)
}
