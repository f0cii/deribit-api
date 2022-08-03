package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/KyberNetwork/deribit-api/pkg/models"
	"github.com/KyberNetwork/deribit-api/pkg/multicast"
	ws "github.com/KyberNetwork/deribit-api/pkg/websocket"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const ()

var (
	debug              = flag.Bool("debug", true, "Enable debug logs")
	wsEndpoint         = flag.String("websocket", "ws://193.58.254.1:8022/ws/api/v2", "Websocket API endpoint")
	apiKey             = flag.String("api-key", "", "API client ID")
	secretKey          = flag.String("secret-key", "", "API secret key")
	ifname             = flag.String("ifname", "bond-colocation", "Interface name to listen for multicast events")
	addrs              = flag.String("addrs", "239.111.111.1,239.111.111.2,239.111.111.3", "UDP addresses to listen to multicast events")
	port               = flag.Int("port", 6100, "UDP port to listen to multicast event")
	gatherDataDuration = flag.Duration("gather-data-duration", 3*time.Minute, "Gather data duration")
	storagePath        = flag.String("storage-path", "examples/", "Path to storage output file")
	log                *zap.SugaredLogger
)

func setupLogger(debug bool) *zap.SugaredLogger {
	pConf := zap.NewProductionEncoderConfig()
	pConf.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(pConf)
	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	if debug {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	l := zap.New(zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level), zap.AddCaller())
	zap.ReplaceGlobals(l)
	return zap.S()
}

func saveData(data interface{}, filename string) error {
	f, err := os.Create(*storagePath + filename)
	if err != nil {
		log.Error("Fail to create file", "filename", filename, "error", err)
		return err
	}

	err = json.NewEncoder(f).Encode(data)
	if err != nil {
		log.Error("Fail to write data to file", "error", err)
		return err
	}

	return nil
}

// listen to multicast orderbook events
func listenToOrderbookEvent(ctx context.Context, m *multicast.Client) {

	orderbookChannels := []string{
		"book.BTC-PERPETUAL",
		"book.BTC-1AUG22-29000-P",
	}
	data := make([]models.OrderBookRawNotification, 0)
	listener := func(e *models.OrderBookRawNotification) {
		data = append(data, *e)
	}
	for _, channel := range orderbookChannels {
		m.On(channel, listener)

	}

	<-ctx.Done()

	for _, channel := range orderbookChannels {
		m.Off(channel, listener)
	}
	saveData(data, "orderbook.json")
}

// listen to multicast trades events
func listenToTradesEvent(ctx context.Context, m *multicast.Client) {
	tradesChannels := []string{
		"trade.option.BTC",
		"trade.future.BTC",
	}
	data := make([]models.TradesNotification, 0)
	listener := func(e *models.TradesNotification) {
		data = append(data, *e)
	}
	for _, channel := range tradesChannels {
		m.On(channel, listener)

	}

	<-ctx.Done()

	for _, channel := range tradesChannels {
		m.Off(channel, listener)
	}

	saveData(data, "trades.json")
}

// listen to multicast ticker events
func listenToTickerEvent(ctx context.Context, m *multicast.Client) {
	tickerChannels := []string{
		"ticker.BTC-PERPETUAL",
		"ticker.BTC-1AUG22-29000-P",
	}
	data := make([]models.TickerNotification, 0)
	listener := func(e *models.TickerNotification) {
		data = append(data, *e)
	}
	for _, channel := range tickerChannels {
		m.On(channel, listener)

	}

	<-ctx.Done()

	for _, channel := range tickerChannels {
		m.Off(channel, listener)
	}
	saveData(data, "ticker.json")
}

func main() {
	flag.Parse()
	wsConfig := &ws.Configuration{
		Addr:          *wsEndpoint,
		ApiKey:        *apiKey,
		SecretKey:     *secretKey,
		AutoReconnect: true,
		DebugMode:     true,
	}
	log = setupLogger(*debug)

	wsClient := ws.New(log, wsConfig)
	err := wsClient.Start()
	if err != nil {
		log.Error("failed to start ws client")
		panic(err)
	}

	udpAddrs := strings.Split(*addrs, ",")
	multicastClient, err := multicast.NewClient(*ifname, udpAddrs, *port, wsClient, []string{"BTC"})
	if err != nil {
		log.Errorw("failed to initiate multicast client", "ifname", ifname, "addrs", addrs)
		panic(err)
	}

	ctx := context.Background()
	notifyCtx, stop := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer stop()

	ctxTimeOut, cancel := context.WithTimeout(notifyCtx, *gatherDataDuration)
	defer cancel()

	err = multicastClient.Start(ctx)
	if err != nil {
		log.Errorw("failed to start multicast client", "ifname", ifname, "addrs", addrs)
		panic(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(3) // for orderbook, trades, ticker notifications

	go func() {
		defer wg.Done()
		listenToOrderbookEvent(ctxTimeOut, multicastClient)
		log.Info("gather multicast orderbook notifications successfully")
	}()

	go func() {
		defer wg.Done()
		listenToTradesEvent(ctxTimeOut, multicastClient)
		log.Info("gather multicast trades notifications successfully")
	}()

	go func() {
		defer wg.Done()
		listenToTickerEvent(ctxTimeOut, multicastClient)
		log.Info("gather multicast ticker notifications successfully")
	}()

	wg.Wait()
}
