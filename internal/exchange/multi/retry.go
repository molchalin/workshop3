package multi

import (
	"context"
	"time"

	"github.com/molchalin/workshop3/internal/exchange"
	"github.com/molchalin/workshop3/pkg/pool"
)

var _ exchange.Client = (*multiClient)(nil)

type multiClient struct {
	clients       map[string]exchange.Client
	pool          *pool.Pool
	retryDuration time.Duration
}

func NewClient(clients map[string]exchange.Client, addrs []string, unavailableDuration, retryDuration time.Duration, permutator pool.Permutator) *multiClient {
	return &multiClient{
		clients:       clients,
		pool:          pool.NewPool(unavailableDuration, addrs, permutator),
		retryDuration: retryDuration,
	}
}

type exchangeResult struct {
	addr string
	rate float64
	err  error
}

func (c *multiClient) ExchangeRate(ctx context.Context, from, to string) (float64, error) {
	ctx, cancel := context.WithCancel(ctx)

	it := c.pool.Iterator()
	timer := time.NewTimer(c.retryDuration)

	var cnt int
	var globalRes *exchangeResult
	resCh := make(chan exchangeResult, len(c.clients))

	handleResult := func(res exchangeResult) (final bool) {
		cnt--
		if res.err == exchange.ErrServerUnavailable {
			c.pool.MarkUnavailable(res.addr)
			return false
		}
		if globalRes == nil {
			cancel()
			globalRes = &res
		}
		return true
	}

LOOP:
	for it.Next() {
		addr := it.Value
		client := c.clients[addr]
		cnt++
		go func() {
			rate, err := client.ExchangeRate(ctx, from, to)
			resCh <- exchangeResult{addr, rate, err}
		}()

		select {
		case res := <-resCh:
			handle := handleResult(res)
			if handle {
				break LOOP
			}
		case <-timer.C:
		}
		timer.Reset(c.retryDuration)
	}
	timer.Stop()

	for ; cnt > 0; cnt-- {
		res := <-resCh
		_ = handleResult(res)
	}

	if globalRes == nil {
		return 0, exchange.ErrServerUnavailable
	}

	return globalRes.rate, globalRes.err
}
