package helloThere

import (
    "reflect"
    "strconv"
    "testing"
    "time"
)

func TestFilter_Integration(t *testing.T) {
    const N = 256

    filter := NewFilter(N<<1, 24*time.Hour)
    defer filter.Close()

    for i := 0; i < N; i++ {
        data := []byte("35.190.247." + strconv.Itoa(i)) // 35.190.247.0/24
        filter.Insert(data)
    }

    for i := N - 1; i >= 0; i-- {
        data := []byte("35.190.247." + strconv.Itoa(i))
        if got := filter.Lookup(data); got != true {
            t.Errorf("filter do not contain data inserted: %s", string(data))
        }
    }
}

func TestFilter_RetentionPolicy(t *testing.T) {
    filter := NewFilter(1, 1 * time.Millisecond)
    defer filter.Close()

    data := []byte("1.1.1.1")
    filter.Insert(data)
    time.Sleep(1 * time.Second)

    if filter.Lookup(data) {
        t.Error("no retention policy compliance")
    }
}

type cuckooMock struct {
    data []byte
}

func (c *cuckooMock) Insert(data []byte) bool {
    c.data = data
    return false
}

func (c *cuckooMock) Lookup(data []byte) bool {
    c.data = data
    return false
}

func (c *cuckooMock) Reset() {}

func TestFilter_Insert(t *testing.T) {
    filter := &UserFilter{
        ops: make(chan func(Filter), 1),
        done: make(chan interface{}, 1),
    }
    defer filter.Close()

    data := []byte("1.1.1.1")
    filter.Insert(data)

    cm := new(cuckooMock)
    (<-filter.ops)(cm)
    if !reflect.DeepEqual(cm.data, data) {
        t.Error("Insert() not called")
    }
}

func TestFilter_Lookup(t *testing.T) {
    filter := &UserFilter{
        ops: make(chan func(Filter), 1),
        done: make(chan interface{}, 1),
    }
    defer filter.Close()

    data := []byte("1.1.1.1")
    go filter.Lookup(data)

    cm := new(cuckooMock)
    (<-filter.ops)(cm)
    if !reflect.DeepEqual(cm.data, data) {
        t.Error("Lookup() not called")
    }
}

func TestFilter_Reset(t *testing.T) {
    filter := NewFilter(10, 24*time.Hour)
    defer filter.Close()

    data := []byte("1.1.1.1")

    filter.Insert(data)
    filter.Reset()

    if got := filter.Lookup(data); got != false {
        t.Errorf("Reset() filter was not resetted")
    }
}

func TestFilter_opsLoop(t *testing.T) {
    filter := NewFilter(1, 24*time.Hour)
    defer filter.Close()

    timeout := time.After(1 * time.Second)
    ok := make(chan interface{}, 1)

    filter.ops <- func(_ Filter) {
        close(ok)
    }

    select {
    case <-timeout:
        t.Errorf("loop timeout before executing a queued operation")
    case <-ok:
    }
}

func BenchmarkUserFilter_Lookup(b *testing.B) {
    filter := NewFilter(1000, 24*time.Hour)
    defer filter.Close()

    for i:=byte(0); i<0xFF; i++ {
        filter.Insert([]byte{i, i, i, i})
    }

    for i:=0; i<b.N; i++ {
        filter.Lookup([]byte{1, 1, 1, 1})
    }
}
