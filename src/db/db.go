package db

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"

	lerrors "github.com/webbtech/contact-api/errors"
	"github.com/webbtech/contact-api/model"
)

const (
	tableName = "WT_Contact"
	dteFmt    = time.RFC1123Z
)

type PersistContact interface {
	Record() *lerrors.StdError
}

type DB struct {
	contact *model.ContactRequest
	item    *item
	svc     *dynamodb.DynamoDB
}

type item struct {
	ID         string
	Date       string
	TimeStamp  int64
	Newsletter bool
	Email      string
	FirstName  string
	LastName   string
	Message    string
	Phone      string
	Type       string
}

func Init(contact *model.ContactRequest) (db *DB) {
	return &DB{contact: contact}
}

func (db *DB) Record() (e *lerrors.StdError) {

	e = db.config()
	if e != nil {
		return e
	}

	av, err := dynamodbattribute.MarshalMap(db.item)
	if err != nil {
		return &lerrors.StdError{Caller: "db.Init", Err: err, Msg: "Got error marchalling new item", StatusCode: 500}
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = db.svc.PutItem(input)
	if err != nil {
		return &lerrors.StdError{Caller: "db.Init", Err: err, Msg: "Got error calling PutItem", StatusCode: 500}
	}

	return nil
}

func (db *DB) config() (e *lerrors.StdError) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ca-central-1")})
	if err != nil {
		return &lerrors.StdError{Caller: "db.Init", Err: err, Msg: "Failed to connect to DB", StatusCode: 500}
	}
	db.svc = dynamodb.New(sess)

	contact := db.contact
	t := time.Now()

	var phone string = ""
	if contact.Phone != nil {
		phone = *contact.Phone
	}

	db.item = &item{
		ID:         uuid.New().String(),
		Date:       t.Format(dteFmt),
		TimeStamp:  t.Unix(), // This (should) allow us to use TTL to expire items
		Newsletter: false,    // this field has not been implemented
		Email:      *contact.Email,
		FirstName:  *contact.FirstName,
		LastName:   *contact.LastName,
		Message:    *contact.Message,
		Phone:      phone,
		Type:       *contact.Type,
	}

	return e
}
