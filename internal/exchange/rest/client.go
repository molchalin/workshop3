package rest

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"syscall"

	"github.com/molchalin/workshop3/internal/exchange"
)

var _ exchange.Client = (*restClient)(nil)

type restClient struct {
	httpClient *http.Client
	addr       string
}

func NewClient(httpClient *http.Client, addr string) *restClient {
	return &restClient{
		httpClient: httpClient,
		addr:       addr,
	}
}

func (c *restClient) ExchangeRate(ctx context.Context, from, to string) (float64, error) {
	u := url.URL{
		Scheme: "http",
		Host:   c.addr,
		Path:   path.Join(from, to),
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return 0, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, wrapError(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusInternalServerError {
		return 0, exchange.ErrServerUnavailable
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bad http status: %v", resp.Status)
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(string(buf), 64)
}

func wrapError(err error) error {
	var netErr net.Error
	var scErr syscall.Errno
	if errors.As(err, &netErr) && netErr.Timeout() ||
		errors.As(err, &scErr) && scErr == syscall.ECONNREFUSED {
		return exchange.ErrServerUnavailable
	}
	return err
}
