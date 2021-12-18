package main

import (
	"errors"
	"os"

	lerrors "github.com/webbtech/contact-api/errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/webbtech/contact-api/handlers"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")

	Stage string
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var h handlers.Handler

	switch request.Path {
	case "/contact":
		h = &handlers.Contact{}
	default:
		h = &handlers.Ping{}
	}
	return h.Response(request)
}

func main() {
	lambda.Start(handler)
}

func logError(err *lerrors.StdError) {
	if Stage == "" {
		Stage = os.Getenv("Stage")
	}
	log.Error(err)
}
