package sink

type Sink interface {
	Put(interface{}) error
	PutMany(*[]interface{}) error
	Close() error
}
