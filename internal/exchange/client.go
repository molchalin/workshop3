//go:generate mockgen -destination mock/client_mock.go -package mock . Client
package exchange

import (
	"context"
	"errors"
)

type Client interface {
	ExchangeRate(ctx context.Context, from, to string) (float64, error)
}

var ErrServerUnavailable = errors.New("server is unavailable")
