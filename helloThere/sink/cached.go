package sink

import (
    "log"
    "time"
)

const (
    defaultCacheCapacity = 1000
)

type CachedSink struct {
    ops   chan func(*[]interface{})
    chain Sink
    done  chan interface{}
}

func NewCachedSink(chain Sink, batchTime time.Duration) Sink {
    sink := &CachedSink{
        ops:   make(chan func(*[]interface{})),
        chain: chain,
        done:  make(chan interface{}, 1),
    }

    go sink.processOps()
    go sink.processBatches(batchTime)

    return sink
}

func (sink *CachedSink) processOps() {
    data := make([]interface{}, 0, defaultCacheCapacity)
    for op := range sink.ops {
        op(&data)
    }
}

func (sink *CachedSink) processBatches(batchTime time.Duration) {
    t := time.NewTicker(batchTime)
    for {
        select {
        case <-t.C:
            sink.ops <- func(data *[]interface{}) {
                if len(*data) > 0 {
                    if err := sink.chain.PutMany(data); err != nil {
                        log.Println("error uploading data to chain sink:", err)
                    }
                    *data = make([]interface{}, 0, defaultCacheCapacity)
                }
            }
        case <-sink.done:
            t.Stop()
            return
        }
    }
}

func (sink *CachedSink) Put(d interface{}) error {
    sink.ops <- func(data *[]interface{}) {
        *data = append(*data, d)
    }
    return nil
}

func (sink *CachedSink) PutMany(data *[]interface{}) error {
    return sink.chain.PutMany(data)
}

func (sink *CachedSink) Close() error {
    close(sink.ops)
    close(sink.done)
    return sink.chain.Close()
}
