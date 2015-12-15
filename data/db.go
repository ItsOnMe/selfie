package data

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pressly/selfie/config"
	"upper.io/bond"
	"upper.io/db/postgresql"
)

//Database structire for all connection to db
type Database struct {
	bond.Session

	Sqlx *sql.DB

	User    UserStore
	App     AppStore
	Release ReleaseStore
	Bundle  BundleStore
}

//DB a global refrence to db struct
var (
	DB *Database
)

//BuildDbURL creates a URL postgress URL as string
func BuildDbURL(
	username string,
	password string,
	hosts []string,
	database string) string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s",
		username, password, strings.Join(hosts, ","), database)
}

//NewDB creates a db session
func NewDB(dbURL string) (*Database, error) {
	connURL, err := postgresql.ParseURL(dbURL)

	db := &Database{}
	db.Session, err = bond.Open(postgresql.Adapter, connURL)
	if err != nil {
		return nil, err
	}

	db.Sqlx = db.Session.Driver().(*sql.DB)

	db.User = UserStore{Store: db.Store(`users`)}
	db.App = AppStore{Store: db.Store(`apps`)}
	db.Release = ReleaseStore{Store: db.Store(`releases`)}
	db.Bundle = BundleStore{Store: db.Store(`bundle`)}

	if DB != nil {
		DB.Close()
	}

	DB = db

	return db, nil
}

//NewDBWithConfig calls NewDB with poper URL
func NewDBWithConfig(dbConf *config.Config) (*Database, error) {
	dbURL := BuildDbURL(dbConf.DB.Username, dbConf.DB.Password, dbConf.DB.Hosts, dbConf.DB.Database)
	return NewDB(dbURL)
}
