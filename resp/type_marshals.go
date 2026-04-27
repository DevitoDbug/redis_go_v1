package resp

import (
	"strconv"
)

func (v Value) marshalArrray() []byte {
	response := []byte{}
	response = append(response, ARRAY)
	response = append(response, strconv.Itoa(len(v.Array))...)
	response = append(response, CRLF...)
	for _, value := range v.Array {
		response = append(response, value.Marshal()...)
	}

	return response
}

func (v Value) marshalBulk() []byte {
	response := []byte{}
	response = append(response, BULK)
	response = append(response, strconv.Itoa(len(v.Bulk))...)
	response = append(response, CRLF...)
	response = append(response, v.Bulk...)
	response = append(response, CRLF...)

	return response
}

func (v Value) marshalString() []byte {
	response := []byte{}
	response = append(response, STRING)
	response = append(response, v.Str...)
	response = append(response, CRLF...)
	return response
}

func (v Value) marshalError() []byte {
	response := []byte{}
	response = append(response, ERROR)
	response = append(response, v.Err...)
	response = append(response, CRLF...)

	return response
}

func (v Value) marshalInteger() []byte {
	response := []byte{}
	response = append(response, INTEGER)
	response = strconv.AppendInt(response, int64(v.Num), 10)
	response = append(response, CRLF...)

	return response
}

func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}
