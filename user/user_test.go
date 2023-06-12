package user

import (
	"github.com/asdine/storm/v3"
	"gopkg.in/mgo.v2/bson"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
	os.Remove(dbPath)
}

func cleanDb(b *testing.B) {
	os.Remove(dbPath)
	u := &User{
		ID:   bson.NewObjectId(),
		Name: "John",
		Role: "Tester",
	}
	err := u.Save()
	if err != nil {
		b.Fatalf("Error saving a user: %s", err)
	}
	b.ResetTimer()
}

func BenchmarkCreate(b *testing.B) {
	cleanDb(b)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "John_" + strconv.Itoa(i),
			Role: "Tester",
		}
		b.StartTimer()
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a user: %s", err)
		}
	}
}

func BenchmarkRead(b *testing.B) {
	cleanDb(b)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "John_" + strconv.Itoa(i),
			Role: "Tester",
		}
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a user: %s", err)
		}
		b.StartTimer()
		_, err = One(u.ID)
		if err != nil {
			b.Fatalf("Error retrieving a user: %s", err)
		}
	}
}

func BenchmarkUpdate(b *testing.B) {
	cleanDb(b)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "John_" + strconv.Itoa(i),
			Role: "Tester",
		}
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a user: %s", err)
		}
		b.StartTimer()
		u.Role = "Developer"
		err = u.Save()
		if err != nil {
			b.Fatalf("Error saving a user: %s", err)
		}
	}
}

func BenchmarkDelete(b *testing.B) {
	cleanDb(b)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "John_" + strconv.Itoa(i),
			Role: "Tester",
		}
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a user: %s", err)
		}
		b.StartTimer()
		err = Delete(u.ID)
		if err != nil {
			b.Fatalf("Error deleting a user: %s", err)
		}
	}
}

func BenchmarkCRUD(b *testing.B) {
	os.Remove(dbPath)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "John",
			Role: "Tester",
		}
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a user: %s", err)
		}
		_, err = One(u.ID)
		if err != nil {
			b.Fatalf("Error retrieving a user: %s", err)
		}
		u.Role = "Developer"
		err = u.Save()
		if err != nil {
			b.Fatalf("Error updating a user: %s", err)
		}
		err = Delete(u.ID)
		if err != nil {
			b.Fatalf("Error deleting a user: %s", err)
		}
	}
}

func TestCRUD(t *testing.T) {
	t.Log("Create")
	u := &User{
		ID:   bson.NewObjectId(),
		Name: "John",
		Role: "Tester",
	}
	err := u.Save()
	if err != nil {
		t.Fatalf("Error saving a user: %s", err)
	}

	t.Log("Read")
	u2, err := One(u.ID)
	if err != nil {
		t.Fatalf("Error retrieving a user: %s", err)
	}
	if !reflect.DeepEqual(u, u2) {
		t.Errorf("Expected user to be %#v, got %#v", u, u2)
	}

	t.Log("Update")
	u.Role = "Developer"
	err = u.Save()
	if err != nil {
		t.Fatalf("Error updating a user: %s", err)
	}
	u3, err := One(u.ID)
	if err != nil {
		t.Fatalf("Error retrieving a user: %s", err)
	}
	if !reflect.DeepEqual(u, u3) {
		t.Errorf("Expected user to be %#v, got %#v", u, u3)
	}

	t.Log("Delete")
	err = Delete(u.ID)
	if err != nil {
		t.Fatalf("Error deleting a user: %s", err)
	}
	_, err = One(u.ID)
	if err == nil {
		t.Fatalf("Record should not exist anymore")
	}
	if err != storm.ErrNotFound {
		t.Fatalf("Error retrieving non-existing user: %s", err)
	}

	t.Log("Read All")
	u2.ID = bson.NewObjectId()
	u3.ID = bson.NewObjectId()

	err = u2.Save()
	if err != nil {
		t.Fatalf("Error saving a user: %s", err)
	}

	err = u3.Save()
	if err != nil {
		t.Fatalf("Error saving a user: %s", err)
	}

	users, err := All()
	if err != nil {
		t.Fatalf("Error retrieving all users: %s", err)
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

}
