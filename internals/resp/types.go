package resp

const (
	ARRAY   = '*'
	BULK    = '$'
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
)

const CRLF = "\r\n"

type Value struct {
	Typ   string
	Str   string
	Num   int
	Bulk  string
	Err   string
	Array []Value
}

func (v Value) Marshal() []byte {
	switch v.Typ {
	case "array":
		return v.marshalArrray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "error":
		return v.marshalError()
	case "integer":
		return v.marshalInteger()
	case "null":
		return v.marshalNull()
	default:
		return []byte{}
	}
}
