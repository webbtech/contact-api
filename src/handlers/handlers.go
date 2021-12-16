package handlers

import "github.com/aws/aws-lambda-go/events"

var headers map[string]string = map[string]string{"Content-Type": "application/json"}

type Handler interface {
	Response() (events.APIGatewayProxyResponse, error)
	Process(request events.APIGatewayProxyRequest)
}

type response struct {
	Body       string
	Headers    map[string]string
	StatusCode int
}

type responseBody struct {
	Code    string      `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}
