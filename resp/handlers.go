package resp

import (
	"strings"
)

func pong(v []Value) Value {
	if len(v) <= 1 {
		return Value{Typ: "string", Str: "pong"}
	}

	var outputString strings.Builder
	for _, value := range v[1:] {
		stringRep := string(value.Bulk)

		outputString.WriteString(stringRep + " ")
	}

	return Value{Typ: "string", Str: outputString.String()}
}
