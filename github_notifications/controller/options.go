package controller

import (
	"time"
)

type Option func(*Controller)

func WithFetchInterval(interval time.Duration) Option {
	return func(c *Controller) {
		if interval > 0 {
			c.fetchInterval = interval
		}
	}
}

func WithListenReason(reason ...string) Option {
	return func(c *Controller) { c.listenReasons.InsertSlice(reason) }
}
