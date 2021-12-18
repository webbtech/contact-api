package mail

// inspiration for this found at: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/ses-example-send-email.html

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"

	lerrors "github.com/webbtech/contact-api/errors"
	"github.com/webbtech/contact-api/model"
)

// TODO: create test for this package
type Mail struct {
	contactInfo *model.ContactRequest
}

var (
	Sender    string
	Recipient string
	Subject   string
	HtmlBody  string
	TextBody  string
)

const (
	MaxMsgLength = 500
	// The character encoding for the email.
	CharSet = "UTF-8"
)

func Init(info *model.ContactRequest) *Mail {
	return &Mail{contactInfo: info}
}

func (m *Mail) Send() (res *ses.SendEmailOutput, e *lerrors.StdError) {

	m.setInfo()

	// Check if we have sender and recipient in our environment
	senderEmail, exists := os.LookupEnv("SenderEmail")
	if !exists {
		err := errors.New("Missing environment variable 'SenderEmail'")
		e = &lerrors.StdError{Caller: "mail.Send", Err: err, Msg: "Email failed to send", StatusCode: 500}
		return nil, e
	}
	Sender = fmt.Sprintf("Webbtech Contact <%s>", senderEmail)

	recipientEmail, exists := os.LookupEnv("RecipientEmail")
	if !exists {
		err := errors.New("Missing environment variable 'RecipientEmail'")
		e = &lerrors.StdError{Caller: "mail.Send", Err: err, Msg: "Email failed to send", StatusCode: 500}
		return nil, e
	}
	Recipient = fmt.Sprintf("Webbtech Admin <%s>", recipientEmail)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ca-central-1")})

	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(HtmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(Sender),
	}

	// Attempt to send the email.
	res, err = svc.SendEmail(input)

	if err != nil {
		// cast err to awserr.Error to get the Code and Message from an error.
		if aerr, ok := err.(awserr.Error); ok {
			e = &lerrors.StdError{Caller: "aws.SendEmail", Code: aerr.Code(), Err: aerr, Msg: "Email failed to send", StatusCode: 500}
		} else {
			e = &lerrors.StdError{Caller: "aws.SendEmail", Code: "unknown", Err: err, Msg: "Email failed to send", StatusCode: 500}
		}
	}

	return res, e
}

func (m *Mail) setInfo() {

	info := *m.contactInfo
	name := fmt.Sprintf("%s %s", *info.FirstName, *info.LastName)

	// Ensure we didn't get a larger than allowed message
	msg := *info.Message
	msgLen := len(msg)
	if msgLen > MaxMsgLength {
		msg = msg[:msgLen-MaxMsgLength]
	}

	Subject = fmt.Sprintf("A Contact Request from: %s", name)
	HtmlBody = "<h1>A Contact Request has been made by: " + fmt.Sprintf("%s", name) + "</h1>" +
		"<div>" + fmt.Sprintf("First name: %s", *info.FirstName) + "</div>" +
		"<div>" + fmt.Sprintf("Last name: %s", *info.LastName) + "</div>" +
		"<div>" + fmt.Sprintf("Email: %s", *info.Email) + "</div>" +
		"<div>" + fmt.Sprintf("Phone: %s", *info.Phone) + "</div>" +
		"<div>" + fmt.Sprintf("Request type: %s", *info.Type) + "</div>" +
		"<div><br />" + fmt.Sprintf("Message: <br />%s", strings.Replace(msg, "\n", "<br />", -1)) + "</div>"

	TextBody = fmt.Sprintf("A Contact Request has been made by: %s \n", name) +
		fmt.Sprintf("First name: %s\n", *info.FirstName) +
		fmt.Sprintf("Last name: %s\n", *info.LastName) +
		fmt.Sprintf("Email: %s\n", *info.Email) +
		fmt.Sprintf("Phone: %s\n", *info.Phone) +
		fmt.Sprintf("Request type: %s\n", *info.Type) +
		fmt.Sprintf("\nMessage:\n%s\n", msg)
}
