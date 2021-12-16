package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {

	os.Setenv("SenderEmail", "info@webbtech.io")
	os.Setenv("RecipientEmail", "info@webbtech.io")

	var msg string
	t.Run("Successful ping", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}))
		defer ts.Close()

		r, err := handler(events.APIGatewayProxyRequest{Path: "/"})

		expectedMsg := "Healthy"
		msg = extractMessage(r.Body)
		if msg != expectedMsg {
			t.Fatalf("Expected error message: %s received: %s", expectedMsg, msg)
		}
		if err != nil {
			t.Fatal("Everything should be ok")
		}
	})

	t.Run("Contact request missing request body", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}))
		defer ts.Close()

		r, _ := handler(events.APIGatewayProxyRequest{HTTPMethod: "POST", Path: "/contact", Body: ""})

		expectedMsg := "Missing input values"
		msg = extractMessage(r.Body)
		if msg != expectedMsg {
			t.Fatalf("Expected error message: %s received: %s", expectedMsg, msg)
		}
	})

	t.Run("Contact request missing input values", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}))
		defer ts.Close()

		input := `{"email":"test@me.com", "firstName":"Test","lastName":"Dummy","message":"My simple message"}`
		r, _ := handler(events.APIGatewayProxyRequest{HTTPMethod: "POST", Path: "/contact", Body: input})

		expectedMsg := "Missing type in input"
		msg = extractMessage(r.Body)
		if msg != expectedMsg {
			t.Fatalf("Expected error message: %s received: %s", expectedMsg, msg)
		}
		if r.StatusCode != 400 {
			t.Fatalf("Expected StatusCode to be: %d, recieved: %d", 400, r.StatusCode)
		}
	})

	t.Run("Contact request missing input values", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}))
		defer ts.Close()

		input := `{"email":"test@me.com", "firstName":"Test","lastName":"Dummy","message":"My simple message","phone":"(905) 123-4567","type":"callback"}`
		r, _ := handler(events.APIGatewayProxyRequest{HTTPMethod: "POST", Path: "/contact", Body: input})

		expectedMsg := "Success"
		msg = extractMessage(r.Body)
		if msg != expectedMsg {
			t.Fatalf("Expected error message: %s received: %s", expectedMsg, msg)
		}
		if r.StatusCode != 201 {
			t.Fatalf("Expected StatusCode to be: %d, recieved: %d", 400, r.StatusCode)
		}
	})
}

func TestEnvVars(t *testing.T) {
	os.Setenv("PARAM1", "VALUE12")
	t.Run("Successful ping", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}))
		defer ts.Close()

		_, _ = handler(events.APIGatewayProxyRequest{Path: "/"})
		p, exists := os.LookupEnv("PARAM1")
		if !exists {
			t.Fatalf("Expected value for PARAM1 to be: %s", p)
		}
	})
}

func extractMessage(b string) (msg string) {
	var dat map[string]string
	_ = json.Unmarshal([]byte(b), &dat)
	return dat["message"]
}
