package main

/*
    Developed by "https://github.com/vitorfmc"

    =======================================================
    Overview:
    =======================================================

    This Lambda Function is example of integration with RDS and DynamoDB.
    The idea is: Everytime a table in DynamoDB receive a data, it will send the event
    information to this lambda function, which will persist in a Aurora RDS.

    DynamoDB Stream ==> This Lambda Function ==> AuroraDB

	Obs.: Remember to give SecretsManagerReadWrite, AWSLambdaDynamoDBExecutionRole and
	AmazonRDSDataFullAccess policies to your Lambda Function

 	=======================================================
    Aurora RDS Script:
	=======================================================

	drop table item_history;
	drop table sale_history;

	create table sale_history (
		id varchar(100) NOT NULL PRIMARY KEY,
		sale_id varchar(100) NOT NULL,
		sale_date datetime not null,
		amount float not null
	);

	create table item_history (
		sale_id varchar(100) NOT NULL,
		id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
		code varchar(100) not null,
		quantity integer not null,
		value float not null,
		CONSTRAINT FK_Hist1 FOREIGN KEY (sale_id) REFERENCES sale_history(id)
	);

	insert into sale_history(id, sale_id, sale_date, amount) values ('1', '1', '2008-01-01 00:00:01', 10.2);
	insert into item_history(sale_id,code,quantity,value) values ('1','1',2,5.2);
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
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

	for _, record := range e.Records {
		if record.EventName == "INSERT" {

			fmt.Printf("[INFO] From Aurora: %s\n", record)

			sale, err := createFromMap(record)
			if err != nil {
				fmt.Println(err)
				return
			}

			sqls := make([]string, 0)
			currDate := convertTimeToAuroraDate(sale.SaleDate)
			amount := strconv.FormatFloat(sale.Amount, 'f', 6, 64)
			sqls = append(sqls, "insert into sale_history(id, sale_id, sale_date, amount) values ('"+sale.Id+"','"+sale.SaleId+"','"+currDate+"','"+amount+"');")

			sqlStatement := "insert into item_history(sale_id,code,quantity,value) values "
			for index, item := range *sale.Items {

				value := strconv.FormatFloat(item.Value, 'f', 6, 64)
				qtd := strconv.FormatInt(item.Quantity, 10)
				sqlStatement = sqlStatement + "('" + sale.Id + "','" + item.Code + "','" + qtd + "','" + value + "')"

				if index < len(*sale.Items)-1 {
					sqlStatement = sqlStatement + ","
				}
			}
			sqls = append(sqls, sqlStatement)

			jsonData, err := json.Marshal(sqls)
			fmt.Printf("[INFO] SQL: %+v\n", string(jsonData))

			sendToAurora(sqls)

		} else {
			fmt.Printf("Event: %s\n", record.EventName)
		}
	}

}

func sendToAurora(sqls []string) {
	sess, _ := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})

	transactionId, _ := beginTransaction(sess)

	paramSet := make([][]*rdsdataservice.SqlParameter, 1)
	paramSet[0] = make([]*rdsdataservice.SqlParameter, 1)

	for _, sql := range sqls {
		resp, err1 := rdsdataservice.New(sess).BatchExecuteStatement(&rdsdataservice.BatchExecuteStatementInput{
			Database:      aws.String(os.Getenv("DATABASE_NAME")),
			ResourceArn:   aws.String(os.Getenv("DATABASE_ARN")),
			SecretArn:     aws.String(os.Getenv("ADMIN_USER_SECRET_ARN")),
			Schema:        aws.String(os.Getenv("DATABASE_SCHEMA")),
			Sql:           aws.String(sql),
			TransactionId: transactionId,
			ParameterSets: paramSet,
		})

		if err1 == nil {
			fmt.Println("[INFO] Response:", resp.String())
		} else {
			fmt.Println("[ERROR] sendToAurora:", err1)
		}
	}

	commitTransaction(sess, transactionId)
}

func beginTransaction(sess *session.Session) (*string, error) {
	request, output := rdsdataservice.New(sess).BeginTransactionRequest(&rdsdataservice.BeginTransactionInput{
		Database:    aws.String(os.Getenv("DATABASE_NAME")),
		ResourceArn: aws.String(os.Getenv("DATABASE_ARN")),
		SecretArn:   aws.String(os.Getenv("ADMIN_USER_SECRET_ARN")),
		Schema:      aws.String(os.Getenv("DATABASE_SCHEMA")),
	})

	err := request.Send()
	if err != nil {
		fmt.Println("[ERROR] beginTransaction: ", err)
	}

	return output.TransactionId, err
}

func commitTransaction(sess *session.Session, transactionId *string) error {
	request, _ := rdsdataservice.New(sess).CommitTransactionRequest(&rdsdataservice.CommitTransactionInput{
		ResourceArn:   aws.String(os.Getenv("DATABASE_ARN")),
		SecretArn:     aws.String(os.Getenv("ADMIN_USER_SECRET_ARN")),
		TransactionId: transactionId,
	})

	err := request.Send()
	if err != nil {
		fmt.Println("[ERROR] commitTransaction: ", err)
	}

	return err
}

func createFromMap(record events.DynamoDBEventRecord) (SaleHistory, error) {
	data, _ := json.Marshal(toGenericMap(record.Change.NewImage))
	var result SaleHistory
	err := json.Unmarshal(data, &result)
	return result, err
}

func toGenericMap(record map[string]events.DynamoDBAttributeValue) map[string]interface{} {

	mapInterface := make(map[string]interface{})

	for name, value := range record {

		if value.DataType() == events.DataTypeString {

			/*
				The value DataTypeString from DynamoDB can be a Date or a String
			*/
			tempStr := strings.Split(value.String(), ".")[0]
			layout := "2006-01-02T15:04:05"
			t, err := time.Parse(layout, tempStr)
			if err != nil {
				mapInterface[name] = value.String()
			} else {
				mapInterface[name] = t
			}

		} else if value.DataType() == events.DataTypeNumber {

			/*
				The value DataTypeNumber comming from DynamoDB can be a Float or Integer
			*/
			f, err := strconv.ParseFloat(value.Number(), 64)
			if err != nil {
				i, err := strconv.ParseInt(value.Number(), 10, 64)
				if err == nil {
					mapInterface[name] = i
				}
			} else {
				mapInterface[name] = f
			}

		} else if value.DataType() == events.DataTypeList {
			itensSlice := make([]interface{}, 0)
			for _, singleMap := range value.List() {
				itensSlice = append(itensSlice, toGenericMap(singleMap.Map()))
			}
			mapInterface[name] = itensSlice

		} else if value.DataType() == events.DataTypeMap {
			mapInterface[name] = value.Map()
		}
	}

	fmt.Printf("[INFO] PARSE: %+v\n", mapInterface)

	return mapInterface
}

func convertTimeToAuroraDate(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}
