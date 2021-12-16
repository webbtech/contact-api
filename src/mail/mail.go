package mail

// inspiration for this found at: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/ses-example-send-email.html

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"

	lerrors "github.com/webbtech/contact-api/errors"
)

var Sender string
var Recipient string

/* func init() {
	senderEmail, exists := os.LookupEnv("SenderEmail")
	if !exists {
		log.Fatal("Missing environment variable 'SenderEmail'")
		return
	}
	Sender = senderEmail
} */

const (
	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	// Sender = "info@webbtech.io"
	// Sender = senderEmail

	// Replace recipient@example.com with a "To" address. If your account
	// is still in the sandbox, this address must be verified.
	// Recipient = "rond@webbtech.io"

	// Specify a configuration set. To use a configuration
	// set, comment the next line and line 92.
	//ConfigurationSet = "ConfigSet"

	// The subject line for the email.
	Subject = "Amazon SES Test (AWS SDK for Go)"

	// The HTML body for the email.
	HtmlBody = "<h1>Amazon SES Test Email (AWS SDK for Go)</h1><p>This email was sent with " +
		"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
		"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"

	//The email body for recipients with non-HTML email clients.
	TextBody = "This email was sent with Amazon SES using the AWS SDK for Go."

	// The character encoding for the email.
	CharSet = "UTF-8"
)

func Send() (res *ses.SendEmailOutput, e *lerrors.StdError) {

	// Check if we have sender and recipient in our environment
	senderEmail, exists := os.LookupEnv("SenderEmail")
	if !exists {
		err := errors.New("Missing environment variable 'SenderEmail'")
		e = &lerrors.StdError{Caller: "mail.Send", Err: err, Msg: "Email failed to send", StatusCode: 500}
		return nil, e
	}
	Sender = senderEmail
	recipientEmail, exists := os.LookupEnv("RecipientEmail")
	if !exists {
		err := errors.New("Missing environment variable 'RecipientEmail'")
		e = &lerrors.StdError{Caller: "mail.Send", Err: err, Msg: "Email failed to send", StatusCode: 500}
		return nil, e
	}
	Recipient = recipientEmail

	// Create a new session in the ca-central-1 region.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ca-central-1")})

	// Create an SES session.
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
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
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
	// fmt.Printf("e.Error: %+v\n", e.Error())

	return res, e
}
