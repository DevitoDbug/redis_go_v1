// Package resp - contains all the code related to the serialization and deserialization on the
// buffer
package resp

import (
	"bufio"
	"io"
	"strconv"
)

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{
		reader: bufio.NewReader(rd),
	}
}

func (r *Resp) readLine() ([]byte, int, error) {
	var line []byte

	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}

		if b == '\r' {
			nextChar, err := r.reader.Peek(1)
			if err != nil {
				return nil, 0, err
			}

			if nextChar[0] == '\n' {
				_, err := r.reader.Discard(1)
				if err != nil {
					return nil, 0, err
				}

				break
			}
		}

		line = append(line, b)
	}

	return line, len(line) + 2, nil // +2 because we need to account for CRLF
}

func (r *Resp) readInt() (val int, bytesConsumed int, err error) {
	line, bytesConsumed, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}

	intValue, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return int(intValue), bytesConsumed, nil
}
