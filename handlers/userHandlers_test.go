package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/christianotieno/go-rest-api/user"
	"gopkg.in/mgo.v2/bson"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestBuildUser(t *testing.T) {
	valid := &user.User{
		ID:   bson.NewObjectId(),
		Name: "John",
		Role: "Tester",
	}
	valid2 := &user.User{
		ID:   valid.ID,
		Name: "John",
		Role: "Developer",
	}
	js, err := json.Marshal(valid)
	if err != nil {
		t.Errorf("Error marshalling a valid user: %s", err)
		t.FailNow()
	}
	ts := []struct {
		txt string
		r   *http.Request
		u   *user.User
		err bool
		exp *user.User
	}{
		{
			txt: "nil request",
			err: true,
		},
		{
			txt: "empty request body",
			r:   &http.Request{},
			err: true,
		},
		{
			txt: "empty user",
			r: &http.Request{
				Body: io.NopCloser(bytes.NewBufferString("{}")),
			},
			err: true,
		},
		{
			txt: "malformed data in request body",
			r: &http.Request{
				Body: io.NopCloser(bytes.NewBufferString(`{"id": "1234"}`)),
			},
			u:   &user.User{},
			err: true,
		},
		{
			txt: "valid request body",
			r: &http.Request{
				Body: io.NopCloser(bytes.NewBuffer(js)),
			},
			u:   &user.User{},
			exp: valid,
		},
		{
			txt: "valid partial request",
			r: &http.Request{
				Body: io.NopCloser(bytes.NewBufferString(`{"age": "37", "role": "Developer"}`)),
			},
			u:   valid,
			exp: valid2,
		},
	}

	for _, tc := range ts {
		t.Log(tc.txt)
		err := bodyToUser(tc.r, tc.u)
		if tc.err {
			if err == nil {
				t.Error("Expected an error but did not get one")
			}
			continue
		}

		if err != nil {
			t.Errorf("Did not expect an error but got one: %s", err)
			continue
		}
		if !reflect.DeepEqual(tc.exp, tc.u) {
			t.Errorf("Expected %v but got %v", tc.exp, tc.u)
		}
	}
}
