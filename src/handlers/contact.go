package handlers

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	log "github.com/sirupsen/logrus"

	"github.com/webbtech/contact-api/db"
	lerrors "github.com/webbtech/contact-api/errors"
	"github.com/webbtech/contact-api/mail"
	"github.com/webbtech/contact-api/model"
)

// CallBackKey constant
const CallBackKey = "callback"

// Stage variable
var Stage string

// Contact struct
type Contact struct {
	request  events.APIGatewayProxyRequest
	response events.APIGatewayProxyResponse
	input    *model.ContactRequest
}

// ========================== Public Methods =============================== //

func (c *Contact) process() {

	rb := responseBody{}
	var body []byte
	var stdError *lerrors.StdError
	var statusCode int = 201

	json.Unmarshal([]byte(c.request.Body), &c.input)
	if c.input == nil {
		stdError = &lerrors.StdError{
			Caller:     "handlers.Contact.Process",
			Code:       lerrors.CodeBadInput,
			Err:        errors.New("Missing request body"),
			Msg:        "Missing request body",
			StatusCode: 400,
		}
	}

	// validate input
	if stdError == nil {
		err := c.validateInput()
		if err != nil {
			errors.As(err, &stdError)
		}
	}

	// Send email
	if stdError == nil {
		mail := mail.Init(c.input)
		_, err := mail.Send()
		if err != nil {
			errors.As(err, &stdError)
		}
	}

	// Create db record
	if stdError == nil {
		var r db.PersistContact
		r = db.Init(c.input)
		err := r.Record()
		if err != nil {
			errors.As(err, &stdError)
		}
	}

	// Process error
	if stdError != nil {
		rb.Code = stdError.Code
		rb.Message = stdError.Msg
		statusCode = stdError.StatusCode
		logError(stdError)
	} else {
		rb.Code = "SUCCESS"
		rb.Message = "Success"
	}

	// Now creact the actual response object
	body, _ = json.Marshal(&rb)
	c.response = events.APIGatewayProxyResponse{
		Body:       string(body),
		Headers:    headers,
		StatusCode: statusCode,
	}
}

func (c *Contact) Response(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	c.request = request
	c.process()
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

// NOTE: this could go into it's own package etc. if we decide to use this method
func logError(err *lerrors.StdError) {
	if Stage == "" {
		Stage = os.Getenv("Stage")
	}

	if Stage != "test" {
		log.Error(err)
	}
}
