package handlers

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestValidateNumberErrors(t *testing.T) {

	// var requestBody string

	t.Run("Number of errors", func(t *testing.T) {
		c := &Contact{}
		requestBody := `{"email":"test@me.com", "firstName":"Test","lastName":"Dummy"}`
		json.Unmarshal([]byte(requestBody), &c.input)

		expectedNumErrs := 3
		err := c.validateInput()
		if err == nil {
			t.Fatal("Expected an error")
		}
		numErrs := len(strings.Split(err.Error(), "\n"))

		if expectedNumErrs != numErrs {
			t.Fatalf("Expected number of errors should be: %d, have: %d", expectedNumErrs, numErrs)
		}
	})
}
func TestValidateHavePhone(t *testing.T) {

	t.Run("Wrong type for callback type", func(t *testing.T) {
		c := &Contact{}
		requestBody := `{"email":"test@me.com", "firstName":"Test","lastName":"Dummy","message":"My simple message","type":"callback"}`
		json.Unmarshal([]byte(requestBody), &c.input)

		expectedNumErrs := 1
		err := c.validateInput()
		if err == nil {
			t.Fatal("Expected an error")
		}
		numErrs := len(strings.Split(err.Error(), "\n"))

		if expectedNumErrs != numErrs {
			t.Fatalf("Expected number of errors should be: %d, have: %d", expectedNumErrs, numErrs)
		}
	})
}
