package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

func iPHandler() (events.APIGatewayProxyResponse, error) {
	resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if resp.StatusCode != 200 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(ip) == 0 {
		return events.APIGatewayProxyResponse{}, ErrNoIP
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello, your IP was: %v", string(ip)),
		StatusCode: 200,
	}, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Printf("path: %+v\n", request.Path)

	hdrs := make(map[string]string)
	hdrs["Content-Type"] = "application/json"
	// hdrs["Access-Control-Allow-Origin"] = "*"
	// hdrs["Access-Control-Allow-Methods"] = "GET,OPTIONS,POST,PUT"
	// hdrs["Access-Control-Allow-Headers"] = "Authorization,Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token"

	if request.Path == "/ping" {
		return events.APIGatewayProxyResponse{
			Body:       "pong",
			Headers:    hdrs,
			StatusCode: 200,
		}, nil
	}

	if request.Path == "/error" {
		/* return events.APIGatewayProxyResponse{
			Body:       ErrNon200Response.Error(),
			StatusCode: 500,
		}, ErrNon200Response */
		return events.APIGatewayProxyResponse{
			Body:       "error",
			Headers:    hdrs,
			StatusCode: 404,
		}, nil
	}

	if request.Path == "/hello" {
		return iPHandler()
	}

	return events.APIGatewayProxyResponse{
		Body:       "invalid request",
		StatusCode: 404,
	}, nil

	/* resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if resp.StatusCode != 200 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(ip) == 0 {
		return events.APIGatewayProxyResponse{}, ErrNoIP
	}

	fmt.Printf("request: %+v\n", request.Path)

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello, your IP was: %v", string(ip)),
		StatusCode: 200,
	}, nil */
}

func main() {
	lambda.Start(handler)
}
