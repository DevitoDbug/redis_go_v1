package resp

import (
	"strings"
)

func (r *Resp) loadHandlers() {
	r.Handlers["PING"] = r.pong
	r.Handlers["SET"] = r.set
	r.Handlers["GET"] = r.get
	r.Handlers["HSET"] = r.hSet
	r.Handlers["HGET"] = r.hGet
	r.Handlers["HGETALL"] = r.hGetAll
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

func (r *Resp) hSet(v []Value) Value {
	if len(v) < 3 ||
		v[0].Typ != "bulk" ||
		v[1].Typ != "bulk" ||
		v[2].Typ != "bulk" ||
		strings.TrimSpace(v[0].Bulk) == "" ||
		strings.TrimSpace(v[1].Bulk) == "" ||
		strings.TrimSpace(v[2].Bulk) == "" {
		return Value{Typ: "error", Err: "invalid input, input structure not supported"}
	}

	hKey := v[0].Bulk
	hStringKey := v[1].Bulk
	var hString strings.Builder
	sep := ""

	// Sample input is user name david ochieng oduor
	for _, value := range v[2:] {
		if value.Typ == "bulk" && strings.TrimSpace(value.Bulk) != "" {
			hString.WriteString(sep)
			hString.WriteString(value.Bulk)
			sep = " " // Enable space between words except for the first iteration
		}
	}

	r.storage.HStoreVal(hKey, hStringKey, hString.String())

	r.storage.PrintStore() // for debugging

	return Value{Typ: "string", Str: "Ok"}
}

func (r *Resp) hGet(v []Value) Value {
	if len(v) != 2 ||
		v[0].Typ != "bulk" ||
		v[1].Typ != "bulk" ||
		strings.TrimSpace(v[1].Bulk) == "" ||
		strings.TrimSpace(v[0].Bulk) == "" {
		return Value{Typ: "error", Err: "invalid input, input structure not supported for GET command"}
	}

	hkey := v[0].Bulk
	hStringkey := v[1].Bulk
	value := r.storage.HGetVal(hkey, hStringkey)
	if strings.TrimSpace(value) == "" {
		return Value{Typ: "error", Err: "value not found"}
	}

	return Value{Typ: "string", Str: value}
}

func (r *Resp) hGetAll(v []Value) Value {
	if len(v) != 1 ||
		v[0].Typ != "bulk" ||
		strings.TrimSpace(v[0].Bulk) == "" {
		return Value{Typ: "error", Err: "invalid input, input structure not supported for GET command"}
	}

	hkey := v[0].Bulk
	hValues := r.storage.HGetAllVal(hkey)
	values := []Value{}

	for key, hValue := range hValues {
		values = append(values,
			Value{
				Typ:  "bulk",
				Bulk: key + ": " + hValue,
			},
		)
	}

	return Value{Typ: "array", Array: values}
}
