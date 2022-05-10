package websocket

import (
	"context"
	"strings"

	"github.com/KyberNetwork/deribit-api/pkg/models"
)

func (c *Client) Subscribe(channels []string) error {
	return c.subscribe(channels, true)
}

func (c *Client) subscribe(channels []string, isNewSubscription bool) error {
	l := c.l.With("func", "subscribe")
	var publicChannels []string
	var privateChannels []string

	c.mu.RLock()
	currentMap := c.subscriptionsMap
	c.mu.RUnlock()
	for _, v := range channels {
		if _, ok := currentMap[v]; ok {
			continue
		}
		if strings.HasPrefix(v, "user.") {
			privateChannels = append(privateChannels, v)
		} else {
			publicChannels = append(publicChannels, v)
		}
	}

	if len(publicChannels) > 0 {
		pubSubResp, err := c.publicSubscribe(context.Background(), &models.SubscribeParams{
			Channels: publicChannels,
		})
		if err != nil {
			l.Errorw("error subscribe public", "err", err)
			return err
		}
		c.mu.Lock()
		if isNewSubscription {
			c.subscriptions = append(c.subscriptions, pubSubResp...)
		}
		for _, v := range pubSubResp {
			c.subscriptionsMap[v] = struct{}{}
		}
		c.mu.Unlock()
	}

	if len(privateChannels) > 0 {
		privateSubResp, err := c.privateSubscribe(context.Background(), &models.SubscribeParams{
			Channels: privateChannels,
		})
		if err != nil {
			l.Errorw("error subscribe private", "err", err)
			return err
		}
		c.mu.Lock()
		if isNewSubscription {
			c.subscriptions = append(c.subscriptions, privateSubResp...)
		}
		for _, v := range privateSubResp {
			c.subscriptionsMap[v] = struct{}{}
		}
		c.mu.Unlock()
	}
	return nil
}

func (c *Client) UnSubscribe(channels []string) error {
	l := c.l.With("func", "UnSubscribe")
	var publicChannels []string
	var privateChannels []string

	c.mu.RLock()
	currentMap := c.subscriptionsMap
	c.mu.RUnlock()
	for _, v := range channels {
		if _, ok := currentMap[v]; !ok {
			continue
		}
		if strings.HasPrefix(v, "user.") {
			privateChannels = append(privateChannels, v)
		} else {
			publicChannels = append(publicChannels, v)
		}
	}

	if len(publicChannels) > 0 {
		pubUnsubResp, err := c.publicUnsubscribe(context.Background(), &models.UnsubscribeParams{
			Channels: publicChannels,
		})
		if err != nil {
			l.Errorw("error subscribe public", "err", err)
			return err
		}
		c.mu.Lock()
		for _, v := range pubUnsubResp {
			delete(c.subscriptionsMap, v)
		}
		c.mu.Unlock()
	}

	if len(privateChannels) > 0 {
		privateUnsubResp, err := c.privateUnsubscribe(context.Background(), &models.UnsubscribeParams{
			Channels: privateChannels,
		})
		if err != nil {
			l.Errorw("error subscribe private", "err", err)
			return err
		}
		c.mu.Lock()
		for _, v := range privateUnsubResp {
			delete(c.subscriptionsMap, v)
		}
		c.mu.Unlock()
	}
	if len(publicChannels)+len(privateChannels) > 0 {
		var newSubscriptions []string
		c.mu.Lock()
		for v := range c.subscriptionsMap {
			newSubscriptions = append(newSubscriptions, v)
		}
		c.subscriptions = newSubscriptions
		c.mu.Unlock()
	}
	return nil
}

func (c *Client) publicSubscribe(ctx context.Context, params *models.SubscribeParams) (result models.SubscribeResponse, err error) {
	err = c.Call(ctx, "public/subscribe", params, &result)
	return
}

func (c *Client) publicUnsubscribe(ctx context.Context, params *models.UnsubscribeParams) (result models.UnsubscribeResponse, err error) {
	err = c.Call(ctx, "public/unsubscribe", params, &result)
	return
}

func (c *Client) privateSubscribe(ctx context.Context, params *models.SubscribeParams) (result models.SubscribeResponse, err error) {
	err = c.Call(ctx, "private/subscribe", params, &result)
	return
}

func (c *Client) privateUnsubscribe(ctx context.Context, params *models.UnsubscribeParams) (result models.UnsubscribeResponse, err error) {
	err = c.Call(ctx, "private/unsubscribe", params, &result)
	return
}
