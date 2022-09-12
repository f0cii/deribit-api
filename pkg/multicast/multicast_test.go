package multicast

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"math"
	"net"
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/KyberNetwork/deribit-api/pkg/common"
	"github.com/KyberNetwork/deribit-api/pkg/models"
	"github.com/KyberNetwork/deribit-api/pkg/multicast/sbe"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/ipv4"
)

const (
	instrumentsFilePath = "mock/instruments.json"
)

var errInvalidParam = errors.New("invalid params")

type MockInstrumentsGetter struct{}

func (m *MockInstrumentsGetter) GetInstruments(
	ctx context.Context, params *models.GetInstrumentsParams,
) ([]models.Instrument, error) {
	var allIns, btcIns, ethIns []models.Instrument
	instrumentsBytes, err := ioutil.ReadFile(instrumentsFilePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(instrumentsBytes, &allIns)
	if err != nil {
		return nil, err
	}
	for _, ins := range allIns {
		if ins.BaseCurrency == "BTC" {
			btcIns = append(btcIns, ins)
		} else if ins.BaseCurrency == "ETH" {
			ethIns = append(ethIns, ins)
		}
	}

	switch params.Currency {
	case "BTC":
		return btcIns, nil
	case "ETH":
		return ethIns, nil
	default:
		return nil, errInvalidParam
	}
}

type MulticastTestSuite struct {
	suite.Suite
	m           *sbe.SbeGoMarshaller
	c           *Client
	wrongClient *Client
	ins         []models.Instrument
	insMap      map[uint32]models.Instrument
}

func TestMulticastTestSuite(t *testing.T) {
	suite.Run(t, new(MulticastTestSuite))
}

func (ts *MulticastTestSuite) SetupSuite() {
	require := ts.Require()
	var (
		ifname     = "not-exits-ifname"
		ipAddrs    = []string{"239.111.111.1", "239.111.111.2", "239.111.111.3"}
		port       = 6100
		currencies = []string{"BTC", "ETH"}
	)

	m := sbe.NewSbeGoMarshaller()

	// Error case
	client, err := NewClient(ifname, ipAddrs, port, &MockInstrumentsGetter{}, currencies)
	require.Error(err)
	require.Nil(client)

	// Success case
	ifi, err := net.InterfaceByIndex(1)
	require.NoError(err)
	client, err = NewClient(ifi.Name, ipAddrs, port, &MockInstrumentsGetter{}, currencies)
	require.NoError(err)
	require.NotNil(client)

	wrongClient, err := NewClient("", ipAddrs, port, &MockInstrumentsGetter{}, []string{"SHIB"})
	require.NoError(err)
	require.NotNil(client)

	var allIns []models.Instrument
	instrumentsBytes, err := ioutil.ReadFile(instrumentsFilePath)
	require.NoError(err)

	err = json.Unmarshal(instrumentsBytes, &allIns)
	require.NoError(err)

	insMap := make(map[uint32]models.Instrument)
	for _, ins := range allIns {
		insMap[ins.InstrumentID] = ins
	}

	sort.Slice(allIns, func(i, j int) bool {
		return allIns[i].InstrumentID < allIns[j].InstrumentID
	})

	testSetupConnection(client, require)

	ts.c = client
	ts.wrongClient = wrongClient
	ts.m = m
	ts.ins = allIns
	ts.insMap = insMap
}

func testSetupConnection(c *Client, require *require.Assertions) {
	mu := &sync.RWMutex{}
	expectedIPGroups := []net.IP{
		net.ParseIP("239.111.111.1"), net.ParseIP("239.111.111.2"), net.ParseIP("239.111.111.3"),
	}

	mu.Lock()
	ipGroups, err := c.setupConnection()
	mu.Unlock()
	require.NoError(err)

	mu.Lock()
	require.Equal(ipGroups, expectedIPGroups)
	mu.Unlock()
}

func (ts *MulticastTestSuite) TestGetAllInstruments() {
	require := ts.Require()

	// success case
	ins, err := getAllInstrument(ts.c.instrumentsGetter, ts.c.supportCurrencies)
	require.NoError(err)

	// sort for comparing
	sort.Slice(ins, func(i, j int) bool {
		return ins[i].InstrumentID < ins[j].InstrumentID
	})
	require.Equal(ins, ts.ins)

	// error case
	ins, err = getAllInstrument(ts.c.instrumentsGetter, []string{"SHIB"})
	require.ErrorIs(err, errInvalidParam)
	require.Nil(ins)
}

func (ts *MulticastTestSuite) TestBuildInstrumentsMapping() {
	require := ts.Require()

	// success case
	err := ts.c.buildInstrumentsMapping()
	require.NoError(err)
	require.Equal(ts.c.instrumentsMap, ts.insMap)

	// error case
	err = ts.wrongClient.buildInstrumentsMapping()
	require.ErrorIs(err, errInvalidParam)
}

func (ts *MulticastTestSuite) TestEventEmitter() {
	require := ts.Require()
	event := "Hello world"
	channel := "test.EventEmitter"
	receiveTimes := 0
	consumer := func(s string) {
		receiveTimes++
		require.Equal(s, event)
	}

	ts.c.On(channel, consumer)
	ts.c.Emit(channel, event)
	ts.c.Off(channel, consumer)
	ts.c.Emit(event)

	require.Equal(1, receiveTimes)
}

func (ts *MulticastTestSuite) TestDecodeInstrumentEvent() {
	require := ts.Require()

	instrumentEvent := []byte{
		0x8c, 0x00, 0xe8, 0x03, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x4a, 0x37, 0x03, 0x00,
		0x01, 0x01, 0x00, 0x02, 0x00, 0x05, 0x03, 0x00, 0x45, 0x54, 0x48, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x45, 0x54, 0x48, 0x00, 0x00, 0x00, 0x00, 0x00, 0x55, 0x53, 0x44, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x45, 0x54, 0x48, 0x00, 0x00, 0x00, 0x00, 0x00, 0x45, 0x54, 0x48, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x60, 0x72, 0xf1, 0xba, 0x7f, 0x01, 0x00, 0x00, 0x00, 0x38, 0xae, 0x36, 0x87, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x58, 0xab, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf0, 0x3f,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xf0, 0x3f, 0xfc, 0xa9, 0xf1, 0xd2, 0x4d, 0x62, 0x40, 0x3f,
		0x61, 0x32, 0x55, 0x30, 0x2a, 0xa9, 0x33, 0x3f, 0x61, 0x32, 0x55, 0x30, 0x2a, 0xa9, 0x33, 0x3f,
		0x61, 0x32, 0x55, 0x30, 0x2a, 0xa9, 0x33, 0x3f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x12, 0x45, 0x54, 0x48, 0x2d, 0x33, 0x31, 0x4d,
		0x41, 0x52, 0x32, 0x33, 0x2d, 0x33, 0x35, 0x30, 0x30, 0x2d, 0x50,
	}

	expectedHeader := sbe.MessageHeader{
		BlockLength:      140,
		TemplateId:       1000,
		SchemaId:         1,
		Version:          1,
		NumGroups:        0,
		NumVarDataFields: 1,
	}

	expectOutPut := Event{
		Type: EventTypeInstrument,
		Data: models.Instrument{
			TickSize:             0.0005,
			TakerCommission:      0.0003,
			SettlementPeriod:     "month",
			QuoteCurrency:        "ETH",
			MinTradeAmount:       1,
			MakerCommission:      0.0003,
			Leverage:             0,
			Kind:                 "option",
			IsActive:             true,
			InstrumentID:         210762,
			InstrumentName:       "ETH-31MAR23-3500-P",
			ExpirationTimestamp:  1680249600000,
			CreationTimestamp:    1648108860000,
			ContractSize:         1,
			BaseCurrency:         "ETH",
			BlockTradeCommission: 0.0003,
			OptionType:           "put",
			Strike:               3500,
		},
	}

	bufferData := bytes.NewBuffer(instrumentEvent)

	var header sbe.MessageHeader
	err := header.Decode(ts.m, bufferData)
	require.NoError(err)
	require.Equal(header, expectedHeader)

	events, err := ts.c.decodeInstrumentEvent(ts.m, bufferData, header)
	require.NoError(err)
	require.Equal(events, expectOutPut)
}

// nolint:funlen
func (ts *MulticastTestSuite) TestDecodeOrderbookEvent() {
	require := ts.Require()

	tests := []struct {
		event          []byte
		expectedHeader sbe.MessageHeader
		expectedOutput Event
		expectedError  error
	}{
		{
			[]byte{
				0x1d, 0x00, 0xe9, 0x03, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x96, 0x37, 0x03, 0x00,
				0x77, 0xc4, 0x15, 0x0d, 0x83, 0x01, 0x00, 0x00, 0x3c, 0x25, 0x7a, 0x7f, 0x0b, 0x00, 0x00, 0x00,
				0x3d, 0x25, 0x7a, 0x7f, 0x0b, 0x00, 0x00, 0x00, 0x01, 0x12, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x60, 0x4e, 0xd3, 0x40, 0x00, 0x00, 0x00, 0x00, 0xc0,
				0x4f, 0xed, 0x40,
			},
			sbe.MessageHeader{
				BlockLength:      29,
				TemplateId:       1001,
				SchemaId:         1,
				Version:          1,
				NumGroups:        1,
				NumVarDataFields: 0,
			},
			Event{
				Type: EventTypeOrderBook,
				Data: models.OrderBookRawNotification{
					Timestamp:      1662371873911,
					InstrumentName: "BTC-PERPETUAL",
					PrevChangeID:   49383351612,
					ChangeID:       49383351613,
					Bids: []models.OrderBookNotificationItem{
						{
							Action: "change",
							Price:  19769.5,
							Amount: 60030,
						},
					},
				},
			},
			nil,
		},
		{
			[]byte{
				0x1d, 0x00, 0xe9, 0x03, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x96, 0x37, 0x03, 0x00,
				0x77, 0xc4, 0x15, 0x0d, 0x83, 0x01, 0x00, 0x00, 0x3c, 0x25, 0x7a, 0x7f, 0x0b, 0x00, 0x00, 0x00,
				0x3d, 0x25, 0x7a, 0x7f, 0x0b, 0x00, 0x00, 0x00, 0x01, 0x12, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x60, 0x4e, 0xd3, 0x40, 0x00, 0x00, 0x00, 0x00, 0xc0,
				0x4f, 0xed, 0x40,
			},
			sbe.MessageHeader{
				BlockLength:      29,
				TemplateId:       1001,
				SchemaId:         1,
				Version:          1,
				NumGroups:        1,
				NumVarDataFields: 0,
			},
			Event{
				Type: EventTypeOrderBook,
				Data: models.OrderBookRawNotification{
					Timestamp:      1662371873911,
					InstrumentName: "BTC-PERPETUAL",
					PrevChangeID:   49383351612,
					ChangeID:       49383351613,
					Asks: []models.OrderBookNotificationItem{
						{
							Action: "change",
							Price:  19769.5,
							Amount: 60030,
						},
					},
				},
			},
			nil,
		},
	}

	for _, test := range tests {
		bufferData := bytes.NewBuffer(test.event)

		var header sbe.MessageHeader
		err := header.Decode(ts.m, bufferData)
		require.NoError(err)
		require.Equal(header, test.expectedHeader)

		eventDecoded, err := ts.c.decodeOrderBookEvent(ts.m, bufferData, header)

		require.ErrorIs(err, test.expectedError)
		require.Equal(test.expectedOutput, eventDecoded)
	}
}

func (ts *MulticastTestSuite) TestDecodeTradesEvent() {
	require := ts.Require()

	event := []byte{
		0x04, 0x00, 0xea, 0x03, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x73, 0x7a, 0x03, 0x00,
		0x53, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xfc, 0xa9, 0xf1, 0xd2, 0x4d, 0x62, 0x50,
		0x3f, 0x9a, 0x99, 0x99, 0x99, 0x99, 0x99, 0xc9, 0x3f, 0xad, 0xb3, 0x83, 0x1c, 0x83, 0x01, 0x00,
		0x00, 0x4a, 0x4d, 0xf5, 0x43, 0xf0, 0xe8, 0x54, 0x3f, 0xf6, 0x28, 0x5c, 0x8f, 0x32, 0xb7, 0xd2,
		0x40, 0xda, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xb6, 0x29, 0x9f, 0x0d, 0x00, 0x00, 0x00,
		0x00, 0x03, 0x00, 0x14, 0xae, 0x47, 0xe1, 0x7a, 0x94, 0x4d, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x85,
	}

	expectedHeader := sbe.MessageHeader{
		BlockLength:      4,
		TemplateId:       1002,
		SchemaId:         1,
		Version:          1,
		NumGroups:        1,
		NumVarDataFields: 0,
	}

	expectOutPut := Event{
		Type: EventTypeTrades,
		Data: models.TradesNotification{
			{
				Amount:         0.2,
				BlockTradeID:   "0",
				Direction:      "sell",
				IndexPrice:     19164.79,
				InstrumentName: "BTC-9SEP22-20000-C",
				InstrumentKind: "option",
				IV:             59.16,
				Liquidation:    "none",
				MarkPrice:      0.00127624,
				Price:          0.001,
				TickDirection:  3,
				Timestamp:      1662630736813,
				TradeID:        "228534710",
				TradeSeq:       1498,
			},
		},
	}

	bufferData := bytes.NewBuffer(event)

	var header sbe.MessageHeader
	err := header.Decode(ts.m, bufferData)
	require.NoError(err)
	require.Equal(header, expectedHeader)

	eventDecoded, err := ts.c.decodeTradesEvent(ts.m, bufferData, header)

	require.NoError(err)
	require.Equal(expectOutPut, eventDecoded)
}

// nolint:funlen
func (ts *MulticastTestSuite) TestDecodeTickerEvent() {
	require := ts.Require()

	event := []byte{
		0x85, 0x00, 0xeb, 0x03, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x7a, 0x3c, 0x03, 0x00,
		0x01, 0xc7, 0x59, 0xe5, 0x15, 0x83, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3f,
		0x40, 0x60, 0xe5, 0xd0, 0x22, 0xdb, 0x59, 0x39, 0x40, 0x5e, 0xba, 0x49, 0x0c, 0x02, 0xfb, 0x3a,
		0x40, 0xa8, 0xc6, 0x4b, 0x37, 0x89, 0xa1, 0x25, 0x40, 0x1f, 0x85, 0xeb, 0x51, 0xb8, 0x67, 0x97,
		0x40, 0x4e, 0x62, 0x10, 0x58, 0x39, 0x24, 0x3a, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0x85, 0xeb, 0x51, 0xb8, 0x67, 0x97,
		0x40, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x3d, 0x47, 0xe4, 0xbb, 0x94, 0x6e, 0x37,
		0x40,
	}

	expectedHeader := sbe.MessageHeader{
		BlockLength:      133,
		TemplateId:       1003,
		SchemaId:         1,
		Version:          1,
		NumGroups:        0,
		NumVarDataFields: 0,
	}

	zero := 0.0
	expectOutPut := Event{
		Type: EventTypeTicker,
		Data: models.TickerNotification{
			Timestamp:       1662519695815,
			Stats:           models.Stats{},
			State:           "open",
			SettlementPrice: 23.431957,
			OpenInterest:    31,
			MinPrice:        25.351,
			MaxPrice:        26.9805,
			MarkPrice:       26.1415,
			LastPrice:       10.8155,
			InstrumentName:  "ETH-30SEP22-40000-P",
			IndexPrice:      1497.93,
			Funding8H:       math.NaN(),
			CurrentFunding:  math.NaN(),
			BestBidPrice:    &zero,
			BestBidAmount:   0,
			BestAskPrice:    &zero,
			BestAskAmount:   0,
		},
	}

	bufferData := bytes.NewBuffer(event)

	var header sbe.MessageHeader
	err := header.Decode(ts.m, bufferData)
	require.NoError(err)
	require.Equal(header, expectedHeader)

	eventDecoded, err := ts.c.decodeTickerEvent(ts.m, bufferData, header)
	require.NoError(err)

	// replace NaN value to `0` and pointer to 'nil'
	expectedData := expectOutPut.Data.(models.TickerNotification)
	outputData := eventDecoded.Data.(models.TickerNotification)

	tickerPtr := reflect.TypeOf(&models.TickerNotification{})
	common.ReplaceNaNValueOfStruct(&expectedData, tickerPtr)
	common.ReplaceNaNValueOfStruct(&outputData, tickerPtr)
	expectedData.BestBidPrice = nil
	expectedData.BestAskPrice = nil
	outputData.BestBidPrice = nil
	outputData.BestAskPrice = nil

	require.Equal(expectOutPut.Type, eventDecoded.Type)
	require.Equal(expectedData, outputData)
}

func (ts *MulticastTestSuite) TestDecodeEvent() {
	assert := ts.Assert()

	tests := []struct {
		event       []byte
		expectError error
	}{
		{
			[]byte{
				0x8c, 0x00, 0xe8, 0x03, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00,
			},
			io.EOF, // decodeInstrument
		},
		{
			[]byte{
				0x1d, 0x00, 0xe9, 0x03, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00,
			},
			io.EOF, // decodeOrderbook
		},
		{
			[]byte{
				0x04, 0x00, 0xea, 0x03, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00,
			},
			io.EOF, // decodeTrade
		},
		{
			[]byte{
				0x85, 0x00, 0xeb, 0x03, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			io.EOF, // decodeTicker
		},
		{
			[]byte{
				0x8c, 0x00, 0xe8, 0x04, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00,
			},
			ErrUnsupportedTemplateID,
		},
	}

	for _, test := range tests {
		bufferData := bytes.NewBuffer(test.event)

		var header sbe.MessageHeader
		err := header.Decode(ts.m, bufferData)
		assert.NoError(err)

		decodedEvent, err := ts.c.decodeEvent(ts.m, bufferData, header)
		assert.ErrorIs(err, test.expectError)
		assert.Nil(decodedEvent.Data)
	}
}

// nolint:funlen
func (ts *MulticastTestSuite) TestDecodeEvents() {
	assert := ts.Assert()
	_ = assert

	tests := []struct {
		event          []byte
		expectedOutput []Event
		expectError    error
	}{
		{
			[]byte{
				0x8c, 0x00, 0xe8, 0x03, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00,
			},
			nil,
			io.EOF, // decodeInstrument
		},
		{
			[]byte{
				0x1d, 0x00, 0xe9, 0x03, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00,
			},
			nil,
			io.EOF, // decodeOrderbook
		},
		{
			[]byte{},
			nil,
			nil,
		},
		{
			[]byte{
				0x1d, 0x00, 0xe9, 0x03, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x96, 0x37, 0x03, 0x00,
				0x77, 0xc4, 0x15, 0x0d, 0x83, 0x01, 0x00, 0x00, 0x3c, 0x25, 0x7a, 0x7f, 0x0b, 0x00, 0x00, 0x00,
				0x3d, 0x25, 0x7a, 0x7f, 0x0b, 0x00, 0x00, 0x00, 0x01, 0x12, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x60, 0x4e, 0xd3, 0x40, 0x00, 0x00, 0x00, 0x00, 0xc0,
				0x4f, 0xed, 0x40,
			},
			[]Event{
				{
					Type: EventTypeOrderBook,
					Data: models.OrderBookRawNotification{
						Timestamp:      1662371873911,
						InstrumentName: "BTC-PERPETUAL",
						PrevChangeID:   49383351612,
						ChangeID:       49383351613,
						Bids: []models.OrderBookNotificationItem{
							{
								Action: "change",
								Price:  19769.5,
								Amount: 60030,
							},
						},
					},
				},
			},
			nil,
		},
	}

	for _, test := range tests {
		bufferData := bytes.NewBuffer(test.event)

		decodedEvent, err := ts.c.decodeEvents(ts.m, bufferData)
		assert.ErrorIs(err, test.expectError)
		_ = decodedEvent
		// assert.Nil(decodedEvent.Data)
	}
}

func (ts *MulticastTestSuite) TestReadPackageHeader() {
	require := ts.Require()
	header := []byte{
		0x8c, 0x00, 0xe8, 0x03, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00,
	}
	n, channelID, seq, err := readPackageHeader(bytes.NewBuffer(header))
	require.NoError(err)
	require.Equal(n, uint16(140))
	require.Equal(channelID, uint16(1000))
	require.Equal(seq, uint32(65537))

	n, channelID, seq, err = readPackageHeader(&bytes.Buffer{})
	require.ErrorIs(err, io.EOF)
	require.Equal(n, uint16(0))
	require.Equal(channelID, uint16(0))
	require.Equal(seq, uint32(0))
}

func (ts *MulticastTestSuite) TestHandlePackageHeader() {
	// func (c *Client) handlePackageHeader(reader io.Reader, chanelIDSeq map[uint16]uint32) error {
	tests := []struct {
		event       []byte
		chanelIDSeq map[uint16]uint32
		expectError error
	}{
		{
			[]byte{},
			nil,
			io.EOF,
		},
		{
			[]byte{
				0x1d, 0x00, 0xe9, 0x03, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00,
			},
			map[uint16]uint32{1001: 65537},
			ErrDuplicatedPackage,
		},
		{
			[]byte{
				0x1d, 0x00, 0xe9, 0x03, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
			},
			map[uint16]uint32{1001: 65537},
			ErrConnectionReset,
		},
		{
			[]byte{
				0x1d, 0x00, 0xe9, 0x03, 0x01, 0x00, 0x01, 0x02, 0x01, 0x00, 0x00, 0x00,
			},
			map[uint16]uint32{1001: 65537},
			nil,
		},
	}

	for _, test := range tests {
		err := ts.c.handlePackageHeader(bytes.NewBuffer(test.event), test.chanelIDSeq)
		ts.Assert().ErrorIs(err, test.expectError)
	}
}

// nolint:funlen
func (ts *MulticastTestSuite) TestEmitEvents() {
	events := []Event{
		{
			Type: EventTypeInstrument,
			Data: models.Instrument{
				TickSize:             0.0005,
				TakerCommission:      0.0003,
				SettlementPeriod:     "month",
				QuoteCurrency:        "ETH",
				MinTradeAmount:       1,
				MakerCommission:      0.0003,
				Leverage:             0,
				Kind:                 "option",
				IsActive:             true,
				InstrumentID:         210762,
				InstrumentName:       "ETH-31MAR23-3500-P",
				ExpirationTimestamp:  1680249600000,
				CreationTimestamp:    1648108860000,
				ContractSize:         1,
				BaseCurrency:         "ETH",
				BlockTradeCommission: 0.0003,
				OptionType:           "put",
				Strike:               3500,
			},
		},
		{
			Type: EventTypeOrderBook,
			Data: models.OrderBookRawNotification{
				Timestamp:      1662371873911,
				InstrumentName: "BTC-PERPETUAL",
				PrevChangeID:   49383351612,
				ChangeID:       49383351613,
				Bids: []models.OrderBookNotificationItem{
					{
						Action: "change",
						Price:  19769.5,
						Amount: 60030,
					},
				},
			},
		},
		{
			Type: EventTypeTrades,
			Data: models.TradesNotification{
				{
					Amount:         0.2,
					BlockTradeID:   "0",
					Direction:      "sell",
					IndexPrice:     19164.79,
					InstrumentName: "BTC-9SEP22-20000-C",
					InstrumentKind: "option",
					IV:             59.16,
					Liquidation:    "none",
					MarkPrice:      0.00127624,
					Price:          0.001,
					TickDirection:  3,
					Timestamp:      1662630736813,
					TradeID:        "228534710",
					TradeSeq:       1498,
				},
			},
		},
		{
			Type: EventTypeTicker,
			Data: models.TickerNotification{
				Timestamp:       1662519695815,
				Stats:           models.Stats{},
				State:           "open",
				SettlementPrice: 23.431957,
				OpenInterest:    31,
				MinPrice:        25.351,
				MaxPrice:        26.9805,
				MarkPrice:       26.1415,
				LastPrice:       10.8155,
				InstrumentName:  "ETH-30SEP22-40000-P",
				IndexPrice:      1497.93,
				Funding8H:       math.NaN(),
				CurrentFunding:  math.NaN(),
				BestBidPrice:    nil,
				BestBidAmount:   0,
				BestAskPrice:    nil,
				BestAskAmount:   0,
			},
		},
	}
	ts.c.emitEvents(events)
}

func (ts *MulticastTestSuite) TestHandleUDPPackage() {
	tests := []struct {
		data          []byte
		expectedError error
	}{
		{
			[]byte{},
			io.EOF,
		},
		{
			[]byte{
				0x1d, 0x00, 0xe9, 0x03, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0xb0, 0x3b, 0x03, 0x00,
				0x17, 0xeb, 0x3a, 0x20, 0x83, 0x01, 0x00, 0x00, 0x10, 0xf1, 0x85, 0x86, 0x0b, 0x00, 0x00, 0x00,
				0x26, 0xf1, 0x85, 0x86, 0x0b, 0x00, 0x00, 0x00, 0x01, 0x12, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x01, 0x02, 0x00, 0x00, 0x00, 0x00, 0x80, 0xf2, 0xd2, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa0, 0xf2, 0xd2, 0x40, 0x00, 0x00, 0x00,
				0x00, 0x00, 0xce, 0xd3, 0x40,
			},
			nil,
		},
		{
			[]byte{
				0x1d, 0x00, 0xe9, 0x03, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0xb0, 0x3b, 0x03, 0x00,
				0x17, 0xeb, 0x3a, 0x20, 0x83, 0x01, 0x00, 0x00, 0x10, 0xf1, 0x85, 0x86, 0x0b, 0x00, 0x00, 0x00,
				0x26, 0xf1, 0x85, 0x86, 0x0b, 0x00, 0x00, 0x00, 0x01, 0x12, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x01, 0x02, 0x00, 0x00, 0x00, 0x00, 0x80, 0xf2, 0xd2, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa0, 0xf2, 0xd2, 0x40, 0x00, 0x00, 0x00,
				0x00, 0x00, 0xce, 0xd3, 0x40,
			},
			nil,
		},
	}

	for _, test := range tests {
		err := ts.c.handleUDPPackage(context.Background(), ts.m, map[uint16]uint32{1001: 65537}, test.data)
		ts.Require().ErrorIs(err, test.expectedError)
	}
}

func (ts *MulticastTestSuite) TestHandle() {
	tests := []struct {
		data          []byte
		expectedError error
	}{
		{
			[]byte{
				0x1d, 0x00, 0xe9, 0x03, 0x01, 0x00, 0x01, 0x00, 0x1d, 0x00, 0xe9, 0x03, 0x01, 0x00, 0x01, 0x00,
				0x01, 0x00, 0x00, 0x00, 0xb0, 0x3b, 0x03, 0x00, 0x17, 0xeb, 0x3a, 0x20, 0x83, 0x01, 0x00, 0x00,
				0x10, 0xf1, 0x85, 0x86, 0x0b, 0x00, 0x00, 0x00, 0x26, 0xf1, 0x85, 0x86, 0x0b, 0x00, 0x00, 0x00,
				0x01, 0x12, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x00, 0x00, 0x00, 0x00, 0x80,
				0xf2, 0xd2, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
				0x00, 0xa0, 0xf2, 0xd2, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0xce, 0xd3, 0x40,
			},
			nil,
		},
		{
			[]byte{
				0x1d, 0x00, 0xe9, 0x03, 0x01, 0x00, 0x01, 0x00, 0x1d, 0x00, 0xe9, 0x03, 0x01, 0x00, 0x01, 0x00,
				0x01, 0x00, 0x00, 0x00,
			},
			io.EOF,
		},
	}

	for _, test := range tests {
		err := ts.c.Handle(ts.m, bytes.NewBuffer(test.data), map[uint16]uint32{1000: 65536})
		ts.Require().ErrorIs(err, test.expectedError)
	}
}

func setupIpv4Conn() (*ipv4.PacketConn, error) {
	lc := net.ListenConfig{}
	baseConn, err := lc.ListenPacket(context.Background(), "udp4", "0.0.0.0:3033")
	if err != nil {
		return nil, err
	}

	conn := ipv4.NewPacketConn(baseConn)
	err = conn.SetControlMessage(ipv4.FlagDst, true)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (ts *MulticastTestSuite) TestReadUDPMulticastPackage() {
	require := ts.Require()

	conn, err := setupIpv4Conn()
	require.NoError(err)

	testData := []byte("Hello World!")

	n, err := conn.WriteTo(testData, nil, conn.LocalAddr())
	require.NoError(err)

	data := make([]byte, 1500)
	res, err := readUDPMulticastPackage(conn, nil, data)
	require.Nil(res)
	require.Equal(data[:n], testData)
	require.NoError(err)
}

func (ts *MulticastTestSuite) TestListenToEvents() {
	require := ts.Require()
	mu := &sync.RWMutex{}

	group := net.ParseIP("239.111.111.1")
	dst := &net.UDPAddr{IP: group, Port: ts.c.port}
	err := ts.c.conn.SetMulticastInterface(ts.c.inf)
	require.NoError(err)

	numEvent := 0
	ts.c.On("book.BTC-31MAR23", func(b *models.OrderBookRawNotification) {
		mu.Lock()
		numEvent++
		mu.Unlock()
		require.Equal(b.InstrumentName, "BTC-31MAR23")
	})
	data := []byte{
		0x18, 0x02, 0x03, 0x00, 0x9f, 0x5c, 0xc3, 0x0a,
		0x1d, 0x00, 0xe9, 0x03, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x37, 0x39, 0x03, 0x00,
		0xc3, 0xf7, 0xde, 0x21, 0x83, 0x01, 0x00, 0x00, 0x2e, 0x93, 0x5f, 0x87, 0x0b, 0x00, 0x00, 0x00,
		0x33, 0x93, 0x5f, 0x87, 0x0b, 0x00, 0x00, 0x00, 0x01, 0x12, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x9e, 0xd4, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x10, 0x8d, 0x40, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0xe0, 0x9e, 0xd4, 0x40, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x04, 0xaf, 0x40,
	}

	_, err = ts.c.conn.WriteTo(data, nil, dst)
	require.NoError(err)

	time.Sleep(100 * time.Millisecond)
	mu.Lock()
	require.Equal(numEvent, 1)
	mu.Unlock()
}

func (ts *MulticastTestSuite) TestStartStop() {
	require := ts.Require()

	// wrong client with nil connection
	err := ts.wrongClient.Stop()
	require.Nil(ts.wrongClient.conn)
	require.NoError(err)

	// there is a connection in client
	err = ts.c.Stop()
	require.NotNil(ts.c.conn)
	ts.Require().NoError(err)

	err = ts.c.Start(context.Background())
	require.NoError(err)

	err = ts.wrongClient.Start(context.Background())
	require.ErrorIs(err, errInvalidParam)
}

func (ts *MulticastTestSuite) TestRestartConnection() {
	err := ts.c.restartConnections(context.Background())
	ts.Require().NoError(err)

	err = ts.wrongClient.restartConnections(context.Background())
	ts.Require().ErrorIs(err, errInvalidParam)
}
