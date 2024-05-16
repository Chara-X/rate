package rate

import (
	"context"
	"time"
)

type Limiter struct{ c chan time.Time }

func NewLimiter(rate time.Duration, size int) *Limiter {
	var c = make(chan time.Time, size)
	go func() {
		defer func() {
			recover()
		}()
		for {
			c <- <-time.After(rate)
		}
	}()
	return &Limiter{c}
}
func (lim *Limiter) Tokens() int { return len(lim.c) }
func (lim *Limiter) Allow() bool {
	select {
	case <-lim.c:
		return true
	default:
		return false
	}
}
func (lim *Limiter) Wait(ctx context.Context) error {
	select {
	case <-lim.c:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (lim *Limiter) Close() { close(lim.c) }
