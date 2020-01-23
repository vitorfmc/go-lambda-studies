package main

/*
    Developed by "https://github.com/vitorfmc"
    
    =======================================================
    Overview:
    =======================================================

    This Lambda Function is example of integration with Api Gateway.
    The idea is: Someone or some app make a request to an Api Gateway,
    which in turn calls this function.

    APP <== http/request ==> API Gateway <==> This Lambda Function

    Obs.: Remember to give Api Gateway policies to your Lambda Function
*/

import (
    "encoding/json"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

var movies = []struct {
    ID int `json:"id"`
    Name string `json:"name"`
}{
    {
        ID: 1,
        Name: "Avengers",
    },
    {
        ID: 2,
        Name: "Ant-Man",
    },
    {
        ID: 3,
        Name: "Thor",
    },
    {
        ID: 4,
        Name: "Hulk",
    }, {
        ID: 5,
        Name: "Doctor Strange",
    },
}

func findAll() (events.APIGatewayProxyResponse, error) {
    response, err := json.Marshal(movies)
    if err != nil {
        return events.APIGatewayProxyResponse{}, err
    }

    return events.APIGatewayProxyResponse{
        StatusCode: 200,
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
        Body: string(response),
    }, nil
}

func main() {
    lambda.Start(findAll)
}