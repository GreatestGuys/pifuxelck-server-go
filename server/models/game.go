package models

import (
	"database/sql"
	"math/rand"

	"github.com/GreatestGuys/pifuxelck-server-go/server/db"
	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
	"github.com/GreatestGuys/pifuxelck-server-go/server/models/common"
)

type Game struct {
	ID string `json:"game_id,omitempty"`
}

type NewGame struct {
	Label   string   `json:"label,omitempty"`
	Players []string `json:"players,omitempty"`
}

type NewGameError struct {
	Label   []string `json:"label,omitempty"`
	Players []string `json:"players,omitempty"`
}

func (e NewGameError) Error() string {
	return common.ModelErrorHelper(e)
}

// CreateGame creates a new game where the first turn is a label submitted by
// the given user ID, and the remaining turns are alternating drawing and labels
// with the players corresponding to the entries in the NewGame struct.
func CreateGame(userID string, newGame NewGame) *Errors {
	if newGame.Label == "" {
		log.Debugf("Failed to create game due to lack of label.")
		return &Errors{NewGame: &NewGameError{
			Label: []string{"A label is required to start a game."},
		}}
	}

	if len(newGame.Players) <= 0 {
		log.Debugf("Failed to create game due to lack of players.")
		return &Errors{NewGame: &NewGameError{
			Players: []string{"At least one other player is required."},
		}}
	}

	genericError := []string{"Unable to create a new game at this time."}
	var errors *Errors
	db.WithTx(func(tx *sql.Tx) error {
		res, _ := tx.Exec(
			`INSERT INTO Games (completed_at_id , next_expiration)
			 VALUES (NULL, NOW() + INTERVAL 2 DAY)`)

		gameID, err := res.LastInsertId()
		if err != nil {
			errors = &Errors{App: genericError}
			return err
		}

		// Insert the first turn into the database. This turn will correspond to
		// the label in the new game request and will be logged as being performed
		// by the user that is creating the game.
		_, err = tx.Exec(
			`INSERT INTO Turns
				 (account_id, game_id, is_complete, is_drawing, label, drawing)
				 VALUES (?, ?, 1, 0, ?, '')`,
			userID, gameID, newGame.Label)
		if err != nil {
			errors = &Errors{App: genericError}
			return err
		}

		// Create a turn entry for each player (in a random order) in the Players
		// list of newGame, alternating drawing and label turns.
		order := rand.Perm(len(newGame.Players))
		for i, v := range order {
			playerID := newGame.Players[v]
			isDrawing := i%2 == 0
			_, err := tx.Exec(
				`INSERT INTO Turns
				 ( account_id
				 , game_id
				 , is_complete
				 , is_drawing
				 , label
				 , drawing
				 ) VALUES (?, ?, 0, ?, '', '')`,
				playerID, gameID, isDrawing)
			if err != nil {
				errors = &Errors{NewGame: &NewGameError{
					Players: []string{"No such player id " + playerID + "."},
				}}
				return err
			}
		}

		return nil
	})

	return errors
}

// UpdateGameCompletedAtTime takes a game ID and updates the completion time if
// the game is over, and does nothing otherwise.
func UpdateGameCompletedAtTime(gameID string) *Errors {
	log.Debugf("Checking if game %v needs a completed at id.", gameID)

	var errors *Errors
	errMsg := []string{"Unable to update the game's completion time."}
	db.WithTx(func(tx *sql.Tx) error {
		// This query is a conditional insert that will create an entry in the
		// GamesCompletedAt table if and only if the game with id gameID is
		// complete AND there is not already an entry in GamesCompletedAt for this
		// game.
		res, err := tx.Exec(
			`INSERT INTO GamesCompletedAt (completed_at)
			 (
			    SELECT NOW()
			    FROM Games
			    WHERE (
			        SELECT completed_at_id
			        FROM Games
			        WHERE id = ?) IS NULL
			      AND 1 = (
			           SELECT SUM(is_complete) = COUNT(*)
			           FROM Turns
			           WHERE game_id = ?)
			    LIMIT 1
			 )`,
			gameID, gameID)
		if err != nil {
			log.Warnf("Query to insert completed at id failed, %v.", err)
			errors = &Errors{App: errMsg}
			return err
		}

		// If there is not last inserted ID, then the game was not over. This is
		// fine, just return nil as a success.
		completedAtID, err := res.LastInsertId()
		if err != nil {
			return nil
		}

		// If there IS a completed at id, then update the game to point to this new
		// entry.
		_, err = tx.Exec(
			"UPDATE Games SET completed_at_id = ? WHERE id = ?",
			completedAtID, gameID)
		return err
	})
	return errors
}
