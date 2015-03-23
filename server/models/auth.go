package models

import (
	"database/sql"
	"encoding/base64"

	"github.com/GreatestGuys/pifuxelck-server-go/server/db"
	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
	"github.com/dustin/randbo"
)

// NewAuthToken creates a new authentication token for the given user ID.
// Presenting this token in the x-pifuxelck-auth header will authenticate the
// request as coming from the user with the given id.
func NewAuthToken(id string) (auth string, errors *Errors) {
	pruneAuthTokens()

	log.Debugf("Generating new random token for user with ID %v.", id)
	r := make([]byte, 32)
	_, err := randbo.New().Read(r)
	if err != nil {
		return "", &Errors{App: []string{"Unable to login at this time."}}
	}
	auth = base64.URLEncoding.EncodeToString(r)

	db.WithTx(func(tx *sql.Tx) error {
		_, err := tx.Exec(
			"INSERT INTO Sessions (auth_token, account_id) VALUES (?, ?)", auth, id)

		if err != nil {
			log.Debugf("Unable to create new authentication token, %v.", err)
			errors = &Errors{App: []string{"Unable to login at this time."}}
			return err
		}

		return nil
	})
	return auth, errors
}

// AuthTokenLookup takes an authentication token an returns the user ID that
// corresponds to the given token.
func AuthTokenLookup(auth string) (id string, errors *Errors) {
	pruneAuthTokens()

	db.WithTx(func(tx *sql.Tx) error {
		row := tx.QueryRow(
			"SELECT account_id FROM Sessions WHERE auth_token = ?", auth)

		err := row.Scan(&id)
		if err != nil {
			log.Debugf("Unable to validate authentication token, %v.", err)
			errors = &Errors{App: []string{"Invalid authentication token."}}
			return err
		}

		return nil
	})
	return id, errors
}

func pruneAuthTokens() {
	db.WithDB(func(db *sql.DB) {
		// Prune all existing authentication tokens that are older than 7 days.
		log.Debugf("Pruning all expired authentication tokens.")
		db.Exec("DELETE FROM Sessions WHERE created_at < NOW() - INTERVAL 7 DAY")
	})
}
