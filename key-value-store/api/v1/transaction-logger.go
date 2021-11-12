package apiv1

type EventType byte
const (
	_                     = iota         // iota == 0; ignore the zero value
	EventDelete EventType = iota         // iota == 1
	EventPut                             // iota == 2; implicitly repeat
)

type Event struct {
	Sequence  uint64                // A unique record ID
	EventType EventType             // The action taken
	Key       string                // The key affected by this transaction
	Value     string                // The value of a PUT the transaction
}



type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error
	ReadEvents() (<-chan Event, <-chan error)
	Run()
}
