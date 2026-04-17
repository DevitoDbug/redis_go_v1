package resp

func pong(v []Value) Value {
	return Value{Typ: "string", Str: "pong"}
}
