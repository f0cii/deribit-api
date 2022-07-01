package models

type InstrumentChangeNotification struct {
	InstrumentName string `json:"instrument_name"`
	State          string `json:"state"`
	Timestamp      int64  `json:"timestamp"`
}
