package userservice

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type User struct {
	ID       string `json:"id"`
	FName    string `json:"first_name"`
	MName    string `json:"middle_name"`
	LName    string `json:"last_name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// Value implements Valueer interface to marshal object into []byte type before
// storing into DB.
func (u User) Value() (driver.Value, error) {
	j, err := json.Marshal(u)
	return j, err
}

// Scan implements Scanner interface to Unmarshal return []byte array from DB into User
func (u *User) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	err := json.Unmarshal(source, u)
	if err != nil {
		return err
	}

	return nil
}
