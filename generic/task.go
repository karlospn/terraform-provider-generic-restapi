package generic

import "encoding/json"

type Task struct {
	Id      string
	Payload json.RawMessage
}
