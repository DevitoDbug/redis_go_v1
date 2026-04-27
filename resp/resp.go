// Package resp - contains all the code related to the serialization and deserialization on the
// buffer
package resp

import (
	"bufio"
	"fmt"
	"io"

	"github.com/DevitoDbug/redis_go_v1/storage"
)

type Resp struct {
	reader   *bufio.Reader
	storage  *storage.Storage
	Handlers map[string]func([]Value) Value
}

func NewResp(rd io.Reader, storage *storage.Storage) *Resp {
	resp := &Resp{
		reader:   bufio.NewReader(rd),
		storage:  storage,
		Handlers: make(map[string]func([]Value) Value),
	}
	resp.loadHandlers()

	return resp
}

func (r *Resp) Read() (*Value, error) {
	_type, err := r.reader.ReadByte() // first character should tell the type refer to types.go
	if err != nil {
		return nil, nil
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		return nil, fmt.Errorf("invalid input type")
	}
}

// readArray - reads array type from the input stream
// Example input - *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n
func (r *Resp) readArray() (*Value, error) {
	v := &Value{}
	v.Typ = "array"

	// start has been read
	arraySize, _, err := r.readInt()
	if err != nil {
		return nil, err
	}

	v.Array = make([]Value, 0, arraySize)
	for range arraySize {
		val, err := r.Read()
		if err != nil {
			return nil, err
		}

		if val == nil {
			return nil, fmt.Errorf("missing value in input stream")
		}

		v.Array = append(v.Array, *val)
	}
	return v, nil
}

// readBulk - reads bulk type from the input stream
func (r *Resp) readBulk() (*Value, error) {
	v := &Value{}
	v.Typ = "bulk"

	bulkSize, _, err := r.readInt() // Reads until the \r\n no need to worry about the initial one
	if err != nil {
		return nil, err
	}

	bulk := make([]byte, bulkSize)
	_, err = r.reader.Read(bulk)
	if err != nil {
		return nil, err
	}
	v.Bulk = string(bulk)

	// Read trailing CRLF
	_, _, err = r.readLine()
	if err != nil {
		return nil, err
	}

	return v, nil
}
