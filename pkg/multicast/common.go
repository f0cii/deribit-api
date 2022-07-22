package multicast

func newInstrumentNotificationChannel() string {
	return "instrument"
}

func newOrderBookNotificationChannel(instrument string) string {
	return "book." + instrument
}

func newTradesNotificationChannel(instrument string) string {
	return "trades." + instrument
}

func newTickerNotificationChannel(instrument string) string {
	return "ticker." + instrument
}
