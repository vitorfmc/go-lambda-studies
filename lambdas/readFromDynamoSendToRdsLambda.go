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
	
	drop table nickname;
	drop table client;

	create table client (
		id varchar(100) NOT NULL PRIMARY KEY,
		name varchar(100) not null,
		email varchar(100) not null
	);

	create table nickname (
		id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
		name varchar(100) not null,
		data date not null,
		client_id varchar(100) NOT NULL,
		CONSTRAINT FK_Client1 FOREIGN KEY (client_id) REFERENCES client(id)
	);

	insert into client(id, name,email) values ('234329489012839021','teste','teste');
	insert into nickname(name,data,client_id) values ('teste','2008-01-01 00:00:01','234329489012839021');
*/

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
)

type Client struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Items    *[]Nickname `json:"items"`
}

type Nickname struct {
	Data string `json:"data"`
	Name string `json:"name"`
}

func main() {
	lambda.Start(handleRequest)
}

func sendToAurora(SQLStatement string){
	sess, _ := session.NewSession(&aws.Config{Region: aws.String("us-east-1")},)
	rdsdataservice_client := rdsdataservice.New(sess)

	fmt.Println("statement:", SQLStatement)

	// The SecretArn is generate when you create a secret for db admin user at 'AWS Secrets Manager' 
	req, resp := rdsdataservice_client.ExecuteStatementRequest(&rdsdataservice.ExecuteStatementInput{
		Database:    aws.String("TABLE_NAME"),
		ResourceArn: aws.String("DATABASE_ARN"),
		SecretArn:   aws.String("ADMIN_USER_SECRET_ARN"),
		Sql:         aws.String(SQLStatement),
	})

	err1 := req.Send()
	if err1 == nil { 
		fmt.Println("Response:", resp)
	} else {
		fmt.Println("error:", err1)
	}
}

func handleRequest(ctx context.Context, e events.DynamoDBEvent) {

	for _, record := range e.Records {
		if record.EventName == "INSERT" {
			SQLStatement := "insert into client(id,name,email) values ('" + record.Change.NewImage["id"].String() + "','" + record.Change.NewImage["name"].String() + "','" + record.Change.NewImage["email"].String() + "');"
			
			sendToAurora(SQLStatement)

			for _, value := range record.Change.NewImage["items"].List() {
				nickname := value.Map()
				SQLStatement = "insert into nickname(name,data,client_id) values ('" + nickname["name"].String() + "','" +nickname["data"].String() + "','" + record.Change.NewImage["id"].String() + "');"
				sendToAurora(SQLStatement)
			}

		}else{
			fmt.Printf("Event: %s\n", record.EventName)
		}
	}

}