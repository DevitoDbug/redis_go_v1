package resp

import (
	"strings"
)

func (r *Resp) loadHandlers() {
	r.Handlers["PING"] = r.pong
	r.Handlers["SET"] = r.set
	r.Handlers["GET"] = r.get
	r.Handlers["HSET"] = r.hset
	r.Handlers["HGET"] = r.hget
}

func (r *Resp) pong(v []Value) Value {
	if len(v) == 0 {
		return Value{Typ: "string", Str: "pong"}
	}

	var outputString strings.Builder
	for _, value := range v {
		stringRep := string(value.Bulk)

		outputString.WriteString(stringRep + " ")
	}

	return Value{Typ: "string", Str: outputString.String()}
}

func (r *Resp) set(v []Value) Value {
	if len(v) != 2 ||
		v[0].Typ != "bulk" ||
		v[1].Typ != "bulk" ||
		strings.TrimSpace(v[0].Bulk) == "" ||
		strings.TrimSpace(v[1].Bulk) == "" {
		return Value{Typ: "error", Err: "invalid input, input structure not supported"}
	}

	key := v[0].Bulk
	value := v[1].Bulk
	r.storage.StoreVal(key, value)

	r.storage.PrintStore() // for debugging

	return Value{Typ: "string", Str: "Ok"}
}

func (r *Resp) get(v []Value) Value {
	if len(v) != 1 ||
		v[0].Typ != "bulk" ||
		strings.TrimSpace(v[0].Bulk) == "" {
		return Value{Typ: "error", Err: "invalid input, input structure not supported for GET command"}
	}

	key := v[0].Bulk
	value := r.storage.GetVal(key)
	if strings.TrimSpace(value) == "" {
		return Value{Typ: "error", Err: "value not found"}
	}

	return Value{Typ: "string", Str: value}
}

func (r *Resp) hset(v []Value) Value {
	return Value{}
}

func (r *Resp) hget(v []Value) Value {
	return Value{}
}
