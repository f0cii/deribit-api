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

	newChannels := c.filterChannels(channels, false)
	privateChannels, publicChannels := splitChannels(newChannels)

	if len(publicChannels) > 0 {
		pubSubResp, err := c.publicSubscribe(context.Background(), &models.SubscribeParams{
			Channels: publicChannels,
		})
		if err != nil {
			l.Errorw("error subscribe public", "err", err)
			return err
		}

		if isNewSubscription {
			c.addChannels(pubSubResp)
		}
	}

	if len(privateChannels) > 0 {
		privateSubResp, err := c.privateSubscribe(context.Background(), &models.SubscribeParams{
			Channels: privateChannels,
		})
		if err != nil {
			l.Errorw("error subscribe private", "err", err)
			return err
		}

		if isNewSubscription {
			c.addChannels(privateSubResp)
		}
	}

	return nil
}

func (c *Client) UnSubscribe(channels []string) error {
	l := c.l.With("func", "UnSubscribe")

	oldChannels := c.filterChannels(channels, true)
	privateChannels, publicChannels := splitChannels(oldChannels)

	if len(publicChannels) > 0 {
		pubUnsubResp, err := c.publicUnsubscribe(context.Background(), &models.UnsubscribeParams{
			Channels: publicChannels,
		})
		if err != nil {
			l.Errorw("error subscribe public", "err", err)
			return err
		}

		c.removeChannels(pubUnsubResp)
	}

	if len(privateChannels) > 0 {
		privateUnsubResp, err := c.privateUnsubscribe(context.Background(), &models.UnsubscribeParams{
			Channels: privateChannels,
		})
		if err != nil {
			l.Errorw("error subscribe private", "err", err)
			return err
		}

		c.removeChannels(privateUnsubResp)
	}

	return nil
}

func (c *Client) publicSubscribe(
	ctx context.Context,
	params *models.SubscribeParams,
) (result models.SubscribeResponse, err error) {
	err = c.Call(ctx, "public/subscribe", params, &result)
	return
}

func (c *Client) publicUnsubscribe(
	ctx context.Context,
	params *models.UnsubscribeParams,
) (result models.UnsubscribeResponse, err error) {
	err = c.Call(ctx, "public/unsubscribe", params, &result)
	return
}

func (c *Client) privateSubscribe(
	ctx context.Context,
	params *models.SubscribeParams,
) (result models.SubscribeResponse, err error) {
	err = c.Call(ctx, "private/subscribe", params, &result)
	return
}

func (c *Client) privateUnsubscribe(
	ctx context.Context,
	params *models.UnsubscribeParams,
) (result models.UnsubscribeResponse, err error) {
	err = c.Call(ctx, "private/unsubscribe", params, &result)
	return
}

func (c *Client) filterChannels(channels []string, exists bool) []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	var newChannels []string
	for _, channel := range channels {
		if _, ok := c.subscriptionsMap[channel]; (exists && ok) || (!exists && !ok) {
			newChannels = append(newChannels, channel)
		}
	}

	return newChannels
}

func (c *Client) addChannels(channels []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.subscriptions = append(c.subscriptions, channels...)
	for _, channel := range channels {
		c.subscriptionsMap[channel] = struct{}{}
	}
}

func (c *Client) removeChannels(channels []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, channel := range channels {
		delete(c.subscriptionsMap, channel)
	}

	if len(channels) > 0 {
		subscriptions := make([]string, 0, len(c.subscriptionsMap))
		for channel := range c.subscriptionsMap {
			subscriptions = append(subscriptions, channel)
		}
		c.subscriptions = subscriptions
	}
}

func isPrivateChannel(channel string) bool {
	return strings.HasPrefix(channel, "user.")
}

func splitChannels(channels []string) (privateChannels, publicChannels []string) {
	for _, channel := range channels {
		if isPrivateChannel(channel) {
			privateChannels = append(privateChannels, channel)
		} else {
			publicChannels = append(publicChannels, channel)
		}
	}

	return
}
