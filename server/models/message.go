package models

// Message corresponds to the top level JSON object that is returned by all
// end points.
type Message struct {
	Errors       *Errors      `json:"errors,omitempty"`
	Game         *Game        `json:"game,omitempty"`
	Games        []Game       `json:"games,omitempty"`
	InboxEntries []InboxEntry `json:"inbox_entries,omitempty"`
	InboxEntry   *InboxEntry  `json:"inbox_entry,omitempty"`
	Meta         *Meta        `json:"meta,omitempty"`
	NewGame      *NewGame     `json:"new_game,omitempty"`
	Turn         *Turn        `json:"turn,omitempty"`
	User         *User        `json:"user,omitempty"`
}

// Errors is a union of all possible error types. It is a sub-field of the
// Message type.
type Errors struct {
	App     []string      `json:"application,omitempty"`
	User    *UserError    `json:"user,omitempty"`
	NewGame *NewGameError `json:"new_game,omitempty"`
}
