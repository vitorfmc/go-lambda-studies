package main

/*
    Developed by "https://github.com/vitorfmc"

    =======================================================
    Overview:
    =======================================================

    This Lambda Function is example of integration with  DynamoDB.
    The idea is make a insert into DB.

	Obs.: Remember to give DynamoDB policies to your Lambda Function

*/

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type SaleHistory struct {
	Id       string         `json:"id"`
	SaleId   string         `json:"sale_id"`
	SaleDate time.Time      `json:"sale_date"`
	Amount   float64        `json:"amount"`
	Items    *[]ItemHistory `json:"items"`
}

type ItemHistory struct {
	SaleId   string  `json:"sale_id"`
	Id       string  `json:"id"`
	Code     string  `json:"ean"`
	Quantity int64   `json:"quantity"`
	Value    float64 `json:"value"`
}

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, e events.DynamoDBEvent) {

	sess, err := session.NewSession(&aws.Config{})

	if err != nil {
		fmt.Println("[ERROR]: ", err)
		return
	}

	db := dynamodb.New(sess, aws.NewConfig().WithRegion("us-east-1"))

	currentTime := time.Now().Format("20060102150405")

	saleItens := make([]ItemHistory, 0)

	item := &ItemHistory{Code: "19283129382193822", Quantity: int64(10), Value: float64(10.20)}
	saleItens = append(saleItens, *item)

	item = &ItemHistory{Code: "19283129382193829", Quantity: int64(10), Value: float64(5.21)}
	saleItens = append(saleItens, *item)

	currentId := time.Now().UnixNano() / int64(time.Millisecond)
	history := &SaleHistory{
		Id:       strconv.FormatInt(currentId, 10),
		SaleId:   currentTime,
		SaleDate: time.Now(),
		Amount:   float64(15.41),
		Items:    &saleItens,
	}

	itemMap, err := dynamodbattribute.MarshalMap(history)

	if err != nil {
		fmt.Println("[ERROR]: ", err)
		return
	}

	_, err = db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("ITEM_NAME")),
		Item:      itemMap,
	})
}
