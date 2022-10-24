package pool

import (
	"math/rand"
	"sync"
	"time"
)

type Pool struct {
	lock        sync.Mutex
	unavailable map[string]time.Time

	addrs               []string
	unavailableDuration time.Duration
	permutator          Permutator
}

type Permutator interface {
	Perm(n int) []int
}

type randPermutator struct{}

func (p *randPermutator) Perm(n int) []int {
	return rand.Perm(n)
}

func NewPool(unavailableDuration time.Duration, addrs []string, permutator Permutator) *Pool {
	if permutator == nil {
		permutator = new(randPermutator)
	}
	return &Pool{
		addrs:               addrs,
		unavailable:         make(map[string]time.Time),
		unavailableDuration: unavailableDuration,
		permutator:          permutator,
	}
}

func (p *Pool) available(addr string) bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	if until, ok := p.unavailable[addr]; ok {
		if time.Now().Before(until) {
			return false
		} else {
			delete(p.unavailable, addr)
		}
	}
	return true
}

func (p *Pool) MarkUnavailable(addr string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.unavailable[addr] = time.Now().Add(p.unavailableDuration)
}

func (p *Pool) Iterator() *iterator {
	return &iterator{
		pool: p,
		perm: p.permutator.Perm(len(p.addrs)),
		done: make(map[string]bool),
	}
}

type iterator struct {
	pool *Pool
	perm []int
	done map[string]bool
	i    int

	Value string
}

func (it *iterator) Next() bool {
	for j := 0; j < len(it.perm); j++ {
		it.Value = it.pool.addrs[it.perm[it.i]]
		it.i = (it.i + 1) % len(it.perm)

		if it.pool.available(it.Value) && !it.done[it.Value] {
			it.done[it.Value] = true
			return true
		}
	}
	it.Value = ""
	return false
}
