package handlers

import (
	"github.com/aws/aws-lambda-go/events"
)

func ContactHandler() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "contact handler",
		StatusCode: 200,
	}, nil
}
