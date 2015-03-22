package models

// Message corresponds to the top level JSON object that is returned by all
// end points.
type Message struct {
	Errors *Errors `json:"errors,omitempty"`
	User   *User   `json:"user,omitempty"`
	Meta   *Meta   `json:"meta,omitempty"`
}

// Errors is a union of all possible error types. It is a sub-field of the
// Message type.
type Errors struct {
	User *UserError `json:"user,omitempty"`
	Meta *MetaError `json:"meta,omitempty"`
}
