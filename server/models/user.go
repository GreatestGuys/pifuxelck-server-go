package models

import (
	"github.com/GreatestGuys/pifuxelck-server-go/server/models/common"
)

// User contains all of the identifying information of a pifuxelck player.
type User struct {
	ID          string `json:"id,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Password    string `json:"-"`
}

// UserError is an error type that is returned
type UserError struct {
	ID          []string `json:"id,omitempty"`
	DisplayName []string `json:"display_name,omitempty"`
	Password    []string `json:"password,omitempty"`
}

func (e UserError) Error() string {
	return common.ModelErrorHelper(e)
}

// CreateUser takes a User object and attempts to create a new user with the
// given credentials. This call can fail if the display name is already
// registered, or if the password is not sufficiently complex.
func CreateUser(user User) (*User, *UserError) {
	return nil, &UserError{DisplayName: []string{"Fuck your display name!"}}
}

// UserLogin takes a User object. The display name and the password fields are
// used to authenticate the user, the ID field is ignored.
func UserLogin(user User) (string, *UserError) {
	return "", &UserError{Password: []string{"Fuck your password name!"}}
}
