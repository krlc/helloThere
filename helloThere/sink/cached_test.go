package sink

import (
    "go.uber.org/goleak"
    "reflect"
    "testing"
    "time"
)

type UserData struct {
    city    string
    country string
    geohash string
}

type MockSink struct {
    action chan *UserData
}

func (sink *MockSink) Put(interface{}) error {
    return nil
}

func (sink *MockSink) PutMany(data *[]interface{}) error {
    sink.action <- (*data)[0].(*UserData)
    return nil
}

func (sink *MockSink) Close() error {
    return nil
}

func TestNewCachedSink(t *testing.T) {
    defer goleak.VerifyNone(t)

    mock := &MockSink{
        action: make(chan *UserData, 1),
    }
    defer func() {
        if err := mock.Close(); err != nil {
            t.Fatal(err)
        }
    }()

    sink := NewCachedSink(mock, 100*time.Millisecond)
    defer func() {
        if err := sink.Close(); err != nil {
            t.Fatal(err)
        }
    }()

    data := &UserData{"test", "test", "test"}

    if err := sink.Put(data); err != nil {
        t.Fatal(err)
    }

    if got := <-mock.action; !reflect.DeepEqual(*got, *data) {
        t.Fatalf("NewCachedSink() got = %v, want %v\n", *got, *data)
    }
}

type BenchMockSink struct{}

func (sink *BenchMockSink) Put(interface{}) error {
    return nil
}

func (sink *BenchMockSink) PutMany(_ *[]interface{}) error {
    return nil
}

func (sink *BenchMockSink) Close() error {
    return nil
}

func BenchmarkCachedSink_Put(b *testing.B) {
    mock := new(BenchMockSink)
    sink := NewCachedSink(mock, 100*time.Millisecond)
    defer func() {
        if err := sink.Close(); err != nil {
            b.Fatal(err)
        }
    }()

    var err error

    for i := 0; i < b.N; i++ {
        err = sink.Put(struct{}{})
    }

    if err != nil {
        b.Error(err)
    }
}
