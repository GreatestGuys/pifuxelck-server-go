package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"

	"github.com/GreatestGuys/pifuxelck-server-go/server/log"
	"github.com/prometheus/client_golang/prometheus"

	_ "github.com/go-sql-driver/mysql"
)

// Config defines all the settings that are required to connect to the pifuxelck
// MySQL database.
type Config struct {
	Host     string
	Port     int
	DB       string
	User     string
	Password string
}

var config = (*Config)(nil)
var configOnce sync.Once

var (
	metricTxCommit = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "db_tx_commits",
		Help: "The number committed transactions.",
	})

	metricTxRollback = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "db_tx_rollbacks",
		Help: "The number rolled back transactions.",
	})
)

func init() {
	prometheus.MustRegister(metricTxCommit)
	prometheus.MustRegister(metricTxRollback)
}

// Initialize the database with a configuration. This method must be called
// prior to calling WithDB.
func Init(c Config) {
	configOnce.Do(func() {
		log.Infof("Initializing database.")

		log.Verbosef("Setting the database config as follows:")
		log.Verbosef("{ Host: %v", c.Host)
		log.Verbosef(", Post: %v", c.Port)
		log.Verbosef(", DB:   %v", c.DB)
		log.Verbosef(", User: %v }", c.User)

		config = &c
	})
}

// WithDB takes a function that is immediately invoked with a reference to the
// pifuxelck database. The function should not close the database connection,
// all resource freeing is handled automatically.
//
// If you are performing multiple SQL operations, you likely want to use WithTx
// which will wrap the operations in a transaction and automatically handle
// commits, rollbacks and retries.
func WithDB(f func(*sql.DB)) {
	if config == nil {
		log.Fatalf("WithDB called prior to initialization of the database.")
	}

	connString := config.User
	if config.Password != "" {
		connString = connString + ":" + config.Password
	}
	connString = connString + "@tcp(" + config.Host + ":" + strconv.Itoa(config.Port) + ")/" + config.DB

	// It is important to connect to the database lazily, as the MySQL server is
	// configured to spin down in times of low usage to keep operating costs low.
	con, err := sql.Open("mysql", connString)
	if err != nil {
		log.Errorf("Unable to open a connection to the MySQl server, %v.", err)
		return
	}
	defer con.Close()

	f(con)
}

// WithTX takes a function that is immediately invoked with a reference to a
// transaction on the pifuxelck database. The function should not commit or roll
// back the transaction. If the passed in function returns an error then the
// transaction will be rolled back other wiser it will be committed.
func WithTx(f func(*sql.Tx) error) error {
	if config == nil {
		log.Fatalf("WithTx called prior to initialization of the database.")
	}

	var err error
	WithDB(func(db *sql.DB) {
		var tx *sql.Tx
		tx, err = db.Begin()
		if err != nil {
			return
		}

		defer func() {
			if p := recover(); p != nil {
				switch p := p.(type) {
				case error:
					err = p
				default:
					err = fmt.Errorf("%s", p)
				}
			}
			if err != nil {
				metricTxRollback.Inc()
				tx.Rollback()
				return
			}
			metricTxCommit.Inc()
			err = tx.Commit()
		}()

		err = f(tx)
	})
	return err
}
