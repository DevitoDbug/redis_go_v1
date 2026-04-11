// Package resp - contains all the code related to the serialization and deserialization on the
// buffer
package resp

import (
	"bufio"
	"fmt"
	"io"
)

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{
		reader: bufio.NewReader(rd),
	}
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
	v.typ = "array"

	// start has been read
	arraySize, _, err := r.readInt()
	if err != nil {
		return nil, err
	}

	v.array = make([]Value, arraySize)
	for range arraySize {
		val, err := r.Read()
		if err != nil {
			return nil, err
		}

		if val == nil {
			return nil, fmt.Errorf("missing value in input stream")
		}

		v.array = append(v.array, *val)
	}

	return v, nil
}

// readBulk - reads bulk type from the input stream
func (r *Resp) readBulk() (*Value, error) {
	v := &Value{}
	v.typ = "bulk"

	bulkSize, _, err := r.readInt() // Reads until the \r\n no need to worry about the initial one
	if err != nil {
		return nil, err
	}

	bulk := make([]byte, bulkSize)
	_, err = r.reader.Read(bulk)
	if err != nil {
		return nil, err
	}
	v.bulk = string(bulk)

	// Read trailing CRLF
	_, _, err = r.readLine()
	if err != nil {
		return nil, err
	}

	return v, nil
}
