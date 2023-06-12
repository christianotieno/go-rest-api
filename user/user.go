package user

import (
	"errors"
	"github.com/asdine/storm/v3"
	"gopkg.in/mgo.v2/bson"
)

// User represents a user in the system
type User struct {
	ID   bson.ObjectId `json:"id" storm:"id"`
	Name string        `json:"name"`
	Role string        `json:"role"`
}

const (
	dbPath = "users.db"
)

// Errors used in the applications
var (
	// Returns ErrRecordInvalid when it encounters ivalid record
	ErrRecordInvalid = errors.New("record is invalid")
)

// All retrieves all users from the database
func All() ([]User, error) {
	db, err := storm.Open(dbPath)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	users := []User{}

	err = db.All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// One returns a single user record from the database
func One(id bson.ObjectId) (*User, error) {
	db, err := storm.Open(dbPath)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	user := new(User)

	err = db.One("ID", id, user)

	if err != nil {
		return nil, err
	}
	return user, nil
}

// Delete removes a given user record from the database
func Delete(id bson.ObjectId) error {
	db, err := storm.Open(dbPath)
	if err != nil {
		return err
	}

	defer db.Close()

	user := new(User)

	err = db.One("ID", id, user)
	if err != nil {
		return err
	}
	return db.DeleteStruct(user)
}

// Save updates or creates a given user in the database
func (u *User) Save() error {
	if err := u.validate(); err != nil {
		return err
	}

	db, err := storm.Open(dbPath)
	if err != nil {
		return err
	}

	defer db.Close()

	return db.Save(u)
}

// Validate checks if the user record contains valid data
func (u *User) validate() error {
	if u.Name == "" {
		return ErrRecordInvalid
	}
	return nil
}
