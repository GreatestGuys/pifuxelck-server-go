package models

import (
	"github.com/GreatestGuys/pifuxelck-server-go/server/models/common"
)

// Meta encodes meta data that does not correspond to any particular model.
type Meta struct {
	Auth string `json:"auth,omitempty"`
}

// MetaError encodes errors that do not correspond to any particular model.
type MetaError struct {
	Encoding []string `json:"encoding,omitempty"`
	App      []string `json:"application,omitempty"`
}

func (e MetaError) Error() string {
	return common.ModelErrorHelper(e)
}
