package common

import (
	"encoding/json"
)

// ModelErrorHelper provides an implementation of Error that models can use to
// easily make their corresponding error types implement the Error interface.
func ModelErrorHelper(model interface{}) string {
	b, err := json.Marshal(model)
	if err != nil {
		return "unmarshallable error"
	}
	return string(b)
}
