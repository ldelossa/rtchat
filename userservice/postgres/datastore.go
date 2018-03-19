package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	us "github.com/ldelossa/rtchat/userservice"
	"github.com/lib/pq"
)

const (
	AddUserQuery        = "INSERT INTO users (u) VALUES ($1);"
	UpdateUserQuery     = "UPDATE users SET u = ($1) WHERE (u)->>'id' = ($2);"
	GetUserByIDQuery    = "SELECT u FROM users WHERE (u)->>'id' = ($1);"
	GetUserByUserName   = "SELECT u FROM users WHERE (u)->>'username' = ($1)"
	DeleteUserByIDQuery = "DELETE FROM users WHERE (u)->>'id' = ($1);"
)

// PGDataStore implements the DataStore interface.
// Embeds the sql.DB type to act as a DB instance
// We make use of jsonb fieds
type DataStore struct {
	connString string
	*sql.DB
}

func NewDatastore(connString string) (*DataStore, error) {
	ds := &DataStore{
		connString: connString,
	}

	err := ds.Init()
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func (ds *DataStore) Init() error {
	// call Open to return init'd DB object
	var err error
	ds.DB, err = sql.Open("postgres", ds.connString)
	if err != nil {
		return err
	}

	// Determine if db can be reached.
	err = ds.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (ds *DataStore) AddUser(u us.User) error {
	// Issue insert and handle possible errors
	_, err := ds.Exec(AddUserQuery, u)
	if err != nil {
		if pError, ok := err.(*pq.Error); ok {
			switch pError.Code {
			case "23505":
				errMsg := fmt.Sprintf("field already exists: %s", pError.Constraint)
				return errors.New(errMsg)
			case "23514":
				errMsg := fmt.Sprintf("field is missing: %s", pError.Constraint)
				return errors.New(errMsg)
			}
		}
		return err
	}
	return nil
}

func (ds *DataStore) UpdateUser(u us.User) error {
	// Make sure user has ID
	if u.ID == "" {
		return fmt.Errorf("POSTGRES: attempting to update user without ID: %v", u)
	}

	// Issue update and handle possible errors
	_, err := ds.Exec(UpdateUserQuery, u, u.ID)
	if err != nil {
		if pError, ok := err.(*pq.Error); ok {
			switch pError.Code {
			case "23505":
				errMsg := fmt.Sprintf("field already exists: %s", pError.Constraint)
				return errors.New(errMsg)
			case "23514":
				errMsg := fmt.Sprintf("field is missing: %s", pError.Constraint)
				return errors.New(errMsg)

			}
		}
		return err
	}
	log.Printf("POSTGRES: Successfully updated user ID %s: %s", u.ID, u)

	return nil
}

func (ds *DataStore) GetUserByID(UID string) (*us.User, error) {
	// Create user object
	var u us.User

	// Issue Get query and handle possible errors
	err := ds.QueryRow(GetUserByIDQuery, UID).Scan(&u)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no rows returned for uid %s", UID)
		}
		return nil, err
	}

	return &u, nil
}

func (ds *DataStore) GetUserByUserName(un string) (*us.User, error) {
	// Create user object
	var u us.User

	// Issue Get query and handle possible errors
	err := ds.QueryRow(GetUserByUserName, un).Scan(&u)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no rows returned for username %s", un)
		}
		return nil, err
	}

	return &u, nil

}

func (ds *DataStore) DeleteUserByID(UID string) error {
	// Attempt to issue delete
	res, err := ds.Exec(DeleteUserByIDQuery, UID)
	affectedRows, err2 := res.RowsAffected()
	switch {
	case affectedRows == 0:
		return fmt.Errorf("recipe does not exist in DB, no deletion action taken")
	case err != nil:
		return fmt.Errorf("error returned by ds Exec: %s", err.Error())
	case err2 != nil:
		return fmt.Errorf("error ruturned on RowsAffected() call: %s", err.Error())
	default:
		return nil
	}
}
