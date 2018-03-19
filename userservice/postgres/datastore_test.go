package postgres

import (
	"testing"

	us "github.com/ldelossa/rtchat/userservice"
)

// Integration tests for datastore.go. You will need a live postgres datastore to run these
var ConnString = "user=postgres dbname=userservice password=dev host=localhost sslmode=disable"

var TestUser = us.User{
	ID:       "TestID",
	FName:    "TestFirstName",
	LName:    "TestLastName",
	MName:    "TestMiddleName",
	Email:    "test@test.com",
	Username: "TestUser",
	Password: "TestPassword!",
}

func cleanUpUser(t *testing.T, ds *DataStore) {
	// Cleanup the user
	_, err := ds.Exec("DELETE FROM users WHERE (u)->>'id' = 'TestID';")
	if err != nil {
		t.Fatalf("could not cleanup user - delete manually: %s", err)
	}
}

func TestDataStore(t *testing.T) {
	// Get new datastore
	ds, err := NewDatastore(ConnString)
	if err != nil {
		t.Fatalf("could not create datastore: %s", err)
	}

	// Attempt Init
	err = ds.Init()
	if err != nil {
		t.Fatalf("failed to init datastore: %s", err)
	}

	addUserTest(t, ds)
	updateUserTest(t, ds)
}

func addUserTest(t *testing.T, ds *DataStore) {
	// Create user
	u := TestUser

	// Attempt add of user
	err := ds.AddUser(u)
	if err != nil {
		t.Fatalf("could not add user to database: %s", err)
	}

	// Attempt to add again, should fail
	err = ds.AddUser(u)
	if err == nil {
		t.Fatalf("expected error adding duplicate user but got nil")
	}

	// Cleanup
	cleanUpUser(t, ds)
}

func updateUserTest(t *testing.T, ds *DataStore) {
	// Create user
	u := TestUser

	// Attempt add of user
	err := ds.AddUser(u)
	if err != nil {
		t.Fatalf("could not add user to database: %s", err)
	}

	// Attempt to update user
	u.Email = "updatetest@updatetest.com"
	err = ds.UpdateUser(u)
	if err != nil {
		cleanUpUser(t, ds)
		t.Fatalf("failed updating user: %s", err)
	}

	// confirm user updated
	var uu us.User
	err = ds.QueryRow("SELECT u FROM users WHERE (u)->>'id' = 'TestID';").Scan(&uu)
	if err != nil {
		cleanUpUser(t, ds)
		t.Fatalf("could not get updated user: %s", err)
	}

	if uu.Email != "updatetest@updatetest.com" {
		cleanUpUser(t, ds)
		t.Fatalf("update of user email did not seem to work: %s", uu.Email)
	}

	// Cleanup
	cleanUpUser(t, ds)
}

func getUserByUserIDTest(t *testing.T, ds *DataStore) {
	// Create user
	u := TestUser

	// Add user to database
	_, err := ds.Exec("INSERT INTO users (u) VALUES ($1)", u)
	if err != nil {
		t.Fatalf("could not add user to database: %s", err)
	}

	// Try to get user by ID
	uu, err := ds.GetUserByID(u.ID)
	if err != nil {
		cleanUpUser(t, ds)
		t.Fatal("couldn not retrieve test user from DB: %s", err)
	}

	// Check user is the same
	if uu.ID != u.ID {
		cleanUpUser(t, ds)
		t.Fatal("retrieved user with unespected user ID: %s", uu.ID)
	}

	// Attempt to get nonexistent user id
	uu, err = ds.GetUserByID("nullNull")
	if err == nil {
		cleanUpUser(t, ds)
		t.Fatal("attempt at getting non-exist userid did not return err")
	}

	// Cleanup
	cleanUpUser(t, ds)
}

func getUserByUserNameTest(t *testing.T, ds *DataStore) {
	// Create user
	u := TestUser

	// Add user to database
	_, err := ds.Exec("INSERT INTO users (u) VALUES ($1)", u)
	if err != nil {
		t.Fatalf("could not add user to database: %s", err)
	}

	// Try to get user by ID
	uu, err := ds.GetUserByUserName(u.Username)
	if err != nil {
		cleanUpUser(t, ds)
		t.Fatal("couldn not retrieve test user from DB: %s", err)
	}

	// Check user is the same
	if uu.ID != u.ID {
		cleanUpUser(t, ds)
		t.Fatal("retrieved user with unespected user ID: %s", uu.ID)
	}

	// Attempt to get nonexistent user id
	uu, err = ds.GetUserByUserName("nullNull")
	if err == nil {
		cleanUpUser(t, ds)
		t.Fatal("attempt at getting non-exist userid did not return err")
	}

	// Cleanup
	cleanUpUser(t, ds)

}

func deleteUserTest(t *testing.T, ds *DataStore) {
	// Create user
	u := TestUser

	// Add user to database
	_, err := ds.Exec("INSERT INTO users (u) VALUES ($1)", u)
	if err != nil {
		t.Fatalf("could not add user to database: %s", err)
	}

	// Attempt to delete user
	err = ds.DeleteUserByID(u.ID)
	if err != nil {
		cleanUpUser(t, ds)
		t.Fatalf("recieved error issue delete query: %s", err)
	}

	// Confirm deletions
	var uu us.User
	err = ds.QueryRow("SELECT u FROM users WHERE (u)->>'id' = 'TestID';").Scan(&uu)
	if err == nil {
		cleanUpUser(t, ds)
		t.Fatalf("expecting ErrNoRows returned from query of delet user")
	}

	// No need to cleanup
}
