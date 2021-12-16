package handlers

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	log "github.com/sirupsen/logrus"

	lerrors "github.com/webbtech/contact-api/errors"
	"github.com/webbtech/contact-api/mail"
)

const CallBackKey = "callback"

type ContactRequest struct {
	Email     *string `json:"email"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Message   *string `json:"message"`
	Phone     *string `json:"phone"`
	Type      *string `json:"type"`
}

type Contact struct {
	response events.APIGatewayProxyResponse
	input    *ContactRequest
}

// ========================== Public Methods =============================== //

func (c *Contact) Process(request events.APIGatewayProxyRequest) {

	rb := responseBody{}
	var stdError *lerrors.StdError
	var statusCode int = 201

	json.Unmarshal([]byte(request.Body), &c.input)
	if c.input == nil {
		rb.Message = "Missing input values"
		body, _ := json.Marshal(&rb)
		c.response = events.APIGatewayProxyResponse{
			Body:       string(body),
			Headers:    headers,
			StatusCode: 400,
		}
		return
	}

	// validate input
	err := c.validateInput()
	if err != nil && errors.As(err, &stdError) {
		rb.Message = stdError.Msg
		rb.Code = stdError.Code
		statusCode = stdError.StatusCode
		body, _ := json.Marshal(&rb)

		c.response = events.APIGatewayProxyResponse{
			Body:       string(body),
			Headers:    headers,
			StatusCode: statusCode,
		}

		log.Error(err)
		return
	}

	// send mail
	_, err = mail.Send()
	if err != nil && errors.As(err, &stdError) {
		rb.Message = stdError.Msg
		rb.Code = stdError.Code
		statusCode = stdError.StatusCode
		body, _ := json.Marshal(&rb)

		c.response = events.APIGatewayProxyResponse{
			Body:       string(body),
			Headers:    headers,
			StatusCode: statusCode,
		}

		log.Error(stdError)
		return
	}

	// if we got this far, all is good
	rb.Message = "Success"
	body, _ := json.Marshal(&rb)
	c.response = events.APIGatewayProxyResponse{
		Body:       string(body),
		Headers:    headers,
		StatusCode: statusCode,
	}
}

func (c *Contact) Response() (events.APIGatewayProxyResponse, error) {
	// mail.Send()
	return c.response, nil
}

// ========================== Private Methods ============================== //

func (c *Contact) validateInput() (err *lerrors.StdError) {
	var inputErrs []string

	if c.input.Email == nil {
		inputErrs = append(inputErrs, "Missing email address in input")
	}
	if c.input.FirstName == nil {
		inputErrs = append(inputErrs, "Missing firstName in input")
	}
	if c.input.LastName == nil {
		inputErrs = append(inputErrs, "Missing lastName in input")
	}
	if c.input.Message == nil {
		inputErrs = append(inputErrs, "Missing message in input")
	}
	if c.input.Type == nil {
		inputErrs = append(inputErrs, "Missing type in input")
	}

	// if user requested a callback, ensure we have a phone number
	if c.input.Type != nil && *c.input.Type == CallBackKey && c.input.Phone == nil {
		inputErrs = append(inputErrs, "Phone number required if Callback is requested")
	}

	if len(inputErrs) > 0 {
		error := errors.New(strings.Join(inputErrs, "\n"))
		err = &lerrors.StdError{Caller: "handlers.validateInput", Code: lerrors.CodeBadInput, Err: error, Msg: error.Error(), StatusCode: 400}
		return err
	}

	return nil
}
