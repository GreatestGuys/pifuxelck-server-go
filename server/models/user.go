package models

import (
	"database/sql"
	"fmt"

	"github.com/GreatestGuys/pifuxelck-server-go/server/db"
	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
	"github.com/GreatestGuys/pifuxelck-server-go/server/models/common"
	"golang.org/x/crypto/bcrypt"
)

// User contains all of the identifying information of a pifuxelck player.
type User struct {
	ID          string `json:"id,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Password    string `json:"password,omitempty"`
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
func CreateUser(user User) (_ *User, userErr *UserError) {
	if len(user.Password) < 8 {
		return nil, &UserError{
			Password: []string{"Password must be at least 8 characters."},
		}
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("Unable to hash password, %v.", err)
		return nil, &UserError{Password: []string{"Invalid password."}}
	}

	db.WithTx(func(tx *sql.Tx) error {
		log.Debugf("Request to register the new user %#v.", user.DisplayName)
		res, err := tx.Exec(
			"INSERT INTO Accounts (display_name, password_hash) VALUES (?, ?)",
			user.DisplayName, hash)

		if err != nil {
			log.Debugf("Attempt to re-register the display name %#v.", user.DisplayName)
			userErr = &UserError{DisplayName: []string{"Display name already taken."}}
			return err
		}

		id, err := res.LastInsertId()
		if err != nil {
			log.Warnf("Unable to obtain the last insert ID for new user, %v.", err)
			userErr = &UserError{DisplayName: []string{"Display name already taken."}}
			return err
		}

		user.ID = fmt.Sprintf("%v", id)
		return nil
	})

	user.Password = ""
	return &user, userErr
}

// UserLookupByPassword takes a User object, and returns the ID of the user
// with the matching display name and password.
func UserLookupByPassword(user User) (id string, userErr *UserError) {
	db.WithTx(func(tx *sql.Tx) error {
		log.Debugf("Retrieving password hash for user %#v.", user.DisplayName)
		row := tx.QueryRow(
			"SELECT id, password_hash FROM Accounts WHERE display_name = ?",
			user.DisplayName)

		var hash []byte
		err := row.Scan(&id, &hash)
		if err != nil {
			log.Debugf("Lookup failed, %v.", err.Error())
			userErr = &UserError{DisplayName: []string{"No such user."}}
			return err
		}

		err = bcrypt.CompareHashAndPassword(hash, []byte(user.Password))
		if err != nil {
			log.Debugf("Lookup failed, bad password.")
			userErr = &UserError{Password: []string{"Invalid password."}}
			return err
		}

		return nil
	})
	return id, userErr
}
