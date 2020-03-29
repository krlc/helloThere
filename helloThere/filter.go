package helloThere

import (
	"time"

	cuckoo "github.com/seiflotfy/cuckoofilter"
)

const (
	defaultRetentionPolicy = 24 * time.Hour
)

type Filter interface {
	Insert([]byte) bool
	Lookup([]byte) bool
	Reset()
}

type FilterClosable interface {
	Filter
	Close()
}

type UserFilter struct {
	ops chan func(Filter)
	done chan interface{}
}

func NewFilter(max uint, retentionPolicy time.Duration) *UserFilter {
	uf := &UserFilter{
		ops: make(chan func(Filter)),
		done: make(chan interface{}, 1),
	}

	go uf.processOps(max)
	go uf.processBatches(retentionPolicy)

	return uf
}

func (uf *UserFilter) processOps(max uint) {
	cf := cuckoo.NewFilter(max)
	for op := range uf.ops {
		op(cf)
	}
}

func (uf *UserFilter) processBatches(retentionPolicy time.Duration) {
	t := time.NewTimer(retentionPolicy)
	for {
		select {
		case <-t.C:
			uf.ops <- func(filter Filter) {
				filter.Reset()
			}
		case <-uf.done:
			t.Stop()
			return
		}
	}
}

func (uf *UserFilter) Insert(ip []byte) {
	uf.ops <- func(filter Filter) {
		filter.Insert(ip)
	}
}

func (uf *UserFilter) Lookup(ip []byte) bool {
	result := make(chan bool, 1)

	uf.ops <- func(filter Filter) {
		result <- filter.Lookup(ip)
	}

	return <-result
}

func (uf *UserFilter) Reset() {
	uf.ops <- func(filter Filter) {
		filter.Reset()
	}
}

func (uf *UserFilter) Close() {
	close(uf.done)
	close(uf.ops)
}
