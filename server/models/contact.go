package models

import (
	"database/sql"

	"github.com/GreatestGuys/pifuxelck-server-go/server/db"
	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
)

// ContactLookup looks up a user given a display name.
func ContactLookup(name string) (user *User, userErr *UserError) {
	db.WithDB(func(db *sql.DB) {
		log.Debugf("Looking up user by display name %#v.", name)

		row := db.QueryRow("SELECT id FROM Accounts WHERE display_name = ?", name)

		var id string
		err := row.Scan(&id)
		if err != nil {
			log.Debugf("Unable to find user %#v.", name)
			userErr = &UserError{DisplayName: []string{"No such user."}}
			return
		}

		log.Debugf("Unable to find user %#v with id %v.", name, id)
		user = &User{DisplayName: name, ID: id}
	})

	return user, userErr
}
