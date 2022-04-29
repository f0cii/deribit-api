package models

type GetFundingChartDataParams struct {
	InstrumentName string `json:"instrument_name"`
	Length         string `json:"length,omitempty"`
}
