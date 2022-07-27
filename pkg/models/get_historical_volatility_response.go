package models

import (
	"encoding/json"
	"fmt"
)

type HistoricalVolatility struct {
	Timestamp uint64  `json:"timestamp"`
	Value     float64 `json:"value"`
}

func (h *HistoricalVolatility) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&h.Timestamp, &h.Value}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in Order: %d != %d", g, e)
	}
	return nil
}

type GetHistoricalVolatilityResponse []HistoricalVolatility
