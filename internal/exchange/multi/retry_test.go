package multi

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/molchalin/workshop3/internal/exchange"
	"github.com/molchalin/workshop3/internal/exchange/mock"
)

type testPermutator struct{}

func (p *testPermutator) Perm(n int) (res []int) {
	for i := 0; i < n; i++ {
		res = append(res, i)
	}
	return res
}

func TestSimple(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	from := "EUR"
	to := "USD"
	m := mock.NewMockClient(ctrl)
	m.EXPECT().ExchangeRate(gomock.Any(), from, to).Return(1.11, nil).Times(1)
	address := "dumb"
	client := NewClient(map[string]exchange.Client{address: m}, []string{address}, time.Second, time.Second, new(testPermutator))

	rate, err := client.ExchangeRate(ctx, from, to)
	if err != nil {
		t.Fatalf("exchange err: %v", err)
	}
	if rate != 1.11 {
		t.Errorf("bad rate: %v", rate)
	}
}

func TestSimpleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	from := "EUR"
	to := "USD"
	m := mock.NewMockClient(ctrl)
	wantErr := errors.New("dumb error")
	m.EXPECT().ExchangeRate(gomock.Any(), from, to).Return(0.0, wantErr).Times(1)
	address := "dumb"
	client := NewClient(map[string]exchange.Client{address: m}, []string{address}, time.Second, time.Second, new(testPermutator))

	_, err := client.ExchangeRate(ctx, from, to)
	if err != wantErr {
		t.Fatalf("wrong exchange err: %v", err)
	}
}

func TestOneFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	from := "EUR"
	to := "USD"
	m1 := mock.NewMockClient(ctrl)
	address1 := "dumb1"
	m2 := mock.NewMockClient(ctrl)
	address2 := "dumb2"
	gomock.InOrder(
		m1.EXPECT().ExchangeRate(gomock.Any(), from, to).Return(0.0, exchange.ErrServerUnavailable).Times(1),
		m2.EXPECT().ExchangeRate(gomock.Any(), from, to).Return(1.11, nil).Times(1),
	)
	client := NewClient(map[string]exchange.Client{address1: m1, address2: m2}, []string{address1, address2}, time.Second, 10*time.Millisecond, new(testPermutator))

	rate, err := client.ExchangeRate(ctx, from, to)
	if err != nil {
		t.Fatalf("exchange err: %v", err)
	}
	if rate != 1.11 {
		t.Fatalf("bad rate: %v", rate)
	}
}

func TestOneTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	from := "EUR"
	to := "USD"
	m1 := mock.NewMockClient(ctrl)
	address1 := "dumb1"
	m2 := mock.NewMockClient(ctrl)
	address2 := "dumb2"
	gomock.InOrder(
		m1.EXPECT().ExchangeRate(gomock.Any(), from, to).Return(0.0, exchange.ErrServerUnavailable).Times(1).Do(func(_ context.Context, _, _ string) { time.Sleep(15 * time.Millisecond) }),
		m2.EXPECT().ExchangeRate(gomock.Any(), from, to).Return(1.11, nil).Times(1),
	)
	client := NewClient(map[string]exchange.Client{address1: m1, address2: m2}, []string{address1, address2}, time.Second, 10*time.Millisecond, new(testPermutator))

	rate, err := client.ExchangeRate(ctx, from, to)
	if err != nil {
		t.Fatalf("exchange err: %v", err)
	}
	if rate != 1.11 {
		t.Fatalf("bad rate: %v", rate)
	}
}

func TestTwoTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	from := "EUR"
	to := "USD"
	m1 := mock.NewMockClient(ctrl)
	address1 := "dumb1"
	m2 := mock.NewMockClient(ctrl)
	address2 := "dumb2"
	gomock.InOrder(
		m1.EXPECT().ExchangeRate(gomock.Any(), from, to).Return(1.11, nil).Times(1).Do(func(_ context.Context, _, _ string) { time.Sleep(15 * time.Millisecond) }),
		m2.EXPECT().ExchangeRate(gomock.Any(), from, to).Return(1.13, nil).Times(1).Do(func(_ context.Context, _, _ string) { time.Sleep(15 * time.Millisecond) }),
	)
	client := NewClient(map[string]exchange.Client{address1: m1, address2: m2}, []string{address1, address2}, time.Second, 10*time.Millisecond, new(testPermutator))

	rate, err := client.ExchangeRate(ctx, from, to)
	if err != nil {
		t.Fatalf("exchange err: %v", err)
	}
	if rate != 1.11 {
		t.Fatalf("bad rate: %v", rate)
	}
}

func TestUnavailable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	from := "EUR"
	to := "USD"
	m1 := mock.NewMockClient(ctrl)
	address1 := "dumb1"
	m2 := mock.NewMockClient(ctrl)
	address2 := "dumb2"
	gomock.InOrder(
		m1.EXPECT().ExchangeRate(gomock.Any(), from, to).Return(0.0, exchange.ErrServerUnavailable).Times(1),
		m2.EXPECT().ExchangeRate(gomock.Any(), from, to).Return(1.11, nil).Times(1),

		m2.EXPECT().ExchangeRate(gomock.Any(), from, to).Return(1.13, nil).Times(1),
	)
	client := NewClient(map[string]exchange.Client{address1: m1, address2: m2}, []string{address1, address2}, time.Second, 10*time.Millisecond, new(testPermutator))

	rate, err := client.ExchangeRate(ctx, from, to)
	if err != nil {
		t.Fatalf("exchange err: %v", err)
	}
	if rate != 1.11 {
		t.Fatalf("bad rate: %v", rate)
	}

	rate2, err := client.ExchangeRate(ctx, from, to)
	if err != nil {
		t.Fatalf("exchange err: %v", err)
	}
	if rate2 != 1.13 {
		t.Fatalf("bad rate: %v", rate2)
	}
}
