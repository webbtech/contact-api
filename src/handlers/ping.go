package handlers

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type Ping struct {
	response events.APIGatewayProxyResponse
}

func (c *Ping) Process(request events.APIGatewayProxyRequest) {
	rb := responseBody{Message: "Healthy"}
	body, _ := json.Marshal(&rb)

	c.response = events.APIGatewayProxyResponse{
		Body:       string(body),
		Headers:    headers,
		StatusCode: 200,
	}
}

func (c *Ping) Response() (events.APIGatewayProxyResponse, error) {
	return c.response, nil
}
