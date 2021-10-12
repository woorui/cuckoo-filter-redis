package filter

import (
	"context"
	"sync"
	"time"

	"github.com/linvon/cuckoo-filter"
)

type Filter struct {
	cf    *cuckoo.Filter
	kv    KV
	mu    *sync.RWMutex
	err   error
	every time.Duration
}

func NewFilter(ctx context.Context, kv KV, every time.Duration, initMemFilter *cuckoo.Filter) (*Filter, func(), error) {
	b, err := kv.Get(ctx)
	if err != nil {
		return nil, func() {}, err
	}

	var cf *cuckoo.Filter
	if len(b) != 0 {
		cf, err = cuckoo.Decode(b)
		if err != nil {
			return nil, func() {}, err
		}
	} else {
		cf = initMemFilter
	}

	filter := &Filter{
		err:   nil,
		cf:    cf,
		kv:    kv,
		every: every,
		mu:    &sync.RWMutex{},
	}

	go filter.runFlush(ctx)

	return filter, func() { filter.Flush(ctx) }, nil
}

func (filter *Filter) Add(data []byte) error {
	filter.mu.Lock()
	defer filter.mu.Unlock()
	if filter.err != nil {
		return filter.err
	}
	filter.cf.AddUnique(data)
	return nil
}

func (filter *Filter) Contain(data []byte) (bool, error) {
	filter.mu.RLock()
	defer filter.mu.RUnlock()
	if filter.err != nil {
		return false, filter.err
	}
	return filter.cf.Contain(data), nil
}

func (filter *Filter) Delete(data []byte) (bool, error) {
	filter.mu.RLock()
	defer filter.mu.RUnlock()
	if filter.err != nil {
		return false, filter.err
	}
	return filter.cf.Delete(data), nil
}

func (filter *Filter) Size() uint {
	filter.mu.RLock()
	defer filter.mu.RUnlock()
	return filter.cf.Size()
}

func (filter *Filter) Flush(ctx context.Context) error {
	filter.mu.Lock()
	defer filter.mu.Unlock()
	b, err := filter.cf.Encode()
	if err != nil {
		return err
	}
	return filter.kv.Set(ctx, b)
}

func (filter *Filter) runFlush(ctx context.Context) {
	t := time.NewTicker(filter.every)
	for {
		select {
		case <-ctx.Done():
			t.Stop()
			filter.err = ctx.Err()
		case <-t.C:
			if err := filter.Flush(ctx); err != nil {
				t.Stop()
				filter.err = err
			}
		}
	}
}
