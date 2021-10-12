package filter

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/linvon/cuckoo-filter"
	"github.com/stretchr/testify/assert"
)

type memKV struct {
	value []byte
}

func (kv *memKV) Get(ctx context.Context) ([]byte, error) {
	return kv.value, nil
}

func (kv *memKV) Set(ctx context.Context, value []byte) error {
	kv.value = value
	return nil
}

var testdata, _ = ioutil.ReadFile("./_testdata")

func Test_Filter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	kv := &memKV{}

	filter, close, err := NewFilter(ctx, kv, 2*time.Second, cuckoo.NewFilter(4, 9, 3900, cuckoo.TableTypePacked))
	assert.Equal(t, err, nil)
	defer close()

	a := []byte("A")
	assert.Equal(t, nil, filter.Add(a))

	ok, err := filter.Contain(a)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, uint(1), filter.Size())

	b := []byte("B")
	assert.Equal(t, err, filter.Add(b))

	ok, err = filter.Contain(b)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, uint(2), filter.Size())

	time.Sleep(3 * time.Second)

	assert.Equal(t, testdata, kv.value)

	cancel()
}

func Test_Persistence_Filter(t *testing.T) {
	ctx := context.Background()

	kv := &memKV{value: testdata}

	filter, close, err := NewFilter(ctx, kv, 2*time.Second, cuckoo.NewFilter(4, 9, 3900, cuckoo.TableTypePacked))
	assert.Equal(t, err, nil)
	defer close()

	ok, err := filter.Contain([]byte("A"))
	assert.Equal(t, nil, err)
	assert.Equal(t, true, ok)

	ok, err = filter.Contain([]byte("B"))
	assert.Equal(t, nil, err)
	assert.Equal(t, true, ok)

	assert.Equal(t, uint(2), filter.Size())
}
