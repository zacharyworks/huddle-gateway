package wsockets

import (
	"encoding/json"
	types "github.com/zacharyworks/huddle-shared/data"
)

type action struct {
	Subset  string
	Type    string
	Payload interface{}
}

func newAction(actionSubset string, actionType string, actionPayload interface{}) *action {
	return &action{
		actionSubset,
		actionType,
		actionPayload,
	}
}

func (a action) build() ([]byte, error) {
	action, err := json.Marshal(types.Action{
		a.Subset,
		a.Type,
		a.Payload,
	})

	return action, err
}
