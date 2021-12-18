package db

import (
	"encoding/json"
	"testing"

	"github.com/webbtech/contact-api/model"
)

func TestDB(t *testing.T) {

	testContact := &model.ContactRequest{}
	requestBody := `{"email":"test@me.com", "firstName":"Test","lastName":"Dummy","message":"My simple message","type":"question"}`
	json.Unmarshal([]byte(requestBody), &testContact)

	t.Run("Successfull Init", func(t *testing.T) {
		var db PersistContact
		db = Init(testContact)
		if _, ok := db.(PersistContact); !ok {
			t.Fatalf("Assertion error")
		}
	})

	t.Run("Successful Record Put", func(t *testing.T) {
		db := Init(testContact)
		err := db.Record()
		if err != nil {
			t.Fatalf("Unexpected error %s", err)
		}
	})
}
