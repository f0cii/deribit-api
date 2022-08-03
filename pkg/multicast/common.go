package multicast

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"syscall"
)

const (
	KindAny = "any"
)

var (
	// ErrLostPackage           = errors.New("lost package")
	ErrConnectionReset       = errors.New("connection reset")
	ErrUnsupportedTemplateId = errors.New("unsupported templateId")
	ErrDuplicatedPackage     = errors.New("duplicated package")
	ErrInvalidIpv4Address    = errors.New("invalid ipv4 address")
	ErrOutOfOrder            = errors.New("package out of order")
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

func getCurrencyFromInstrument(instrument string) string {
	return strings.Split(instrument, "-")[0]
}

func isNetConnClosedErr(err error) bool {
	switch {
	case
		errors.Is(err, net.ErrClosed),
		errors.Is(err, io.EOF),
		errors.Is(err, syscall.EPIPE):
		return true
	default:
		return false
	}
}
