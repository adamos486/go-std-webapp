package log

//Entry is a struct for a log entry.
type Entry struct {
	key   string
	value interface{}
}

//EntryInterface is a generic interface for logger entries.
type EntryInterface interface {
	Key() string
	Value() interface{}
}

//NewEntry creates a new struct object of type Entry
func NewEntry(k string, v interface{}) *Entry {
	return &Entry{
		key:   k,
		value: v,
	}
}

//Key fetches the private key from the struct.
func (e *Entry) Key() string {
	return e.key
}

//Value fetches the private value from the struct.
func (e *Entry) Value() interface{} {
	return e.value
}
