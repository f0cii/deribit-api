package websocket

import (
	"testing"

	"github.com/KyberNetwork/deribit-api/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestEventOnEmitOff(t *testing.T) {
	client := newClient()

	channel := "book.BTC-PERPETUAL.100ms"
	eventCh := make(chan *models.OrderBookNotification, 1)
	listener := func(e *models.OrderBookNotification) {
		eventCh <- e
	}
	client.On(channel, listener)

	expect := &models.OrderBookNotification{}
	client.Emit(channel, expect)
	if assert.Len(t, eventCh, 1) {
		event := <-eventCh
		assert.Equal(t, expect, event)
	}

	client.Off(channel, listener)
	client.Emit(expect)
	assert.Len(t, eventCh, 0)
}
