package multicast

import (
	"errors"
	"fmt"
)

const (
	KindAny = "any"
)

var (
	ErrLostPackage           = errors.New("lost package")
	ErrConnectionReset       = errors.New("connection reset")
	ErrUnsupportedTemplateId = errors.New("unsupported templateId")
)

func newInstrumentNotificationChannel(kind, currency string) string {
	return fmt.Sprintf("instrument.%s.%s", kind, currency)
}

func newOrderBookNotificationChannel(instrument string) string {
	return "book." + instrument
}

func newTradesNotificationChannel(kind, currency string) string {
	return fmt.Sprintf("trades.%s.%s", kind, currency)
}

func newTickerNotificationChannel(instrument string) string {
	return "ticker." + instrument
}
