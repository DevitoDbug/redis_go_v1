package resp

var Handler map[string]func([]Value) Value

func init() {
	Handler = make(map[string]func([]Value) Value)
	Handler["PING"] = pong
}
