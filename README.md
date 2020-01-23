# GO REST API Example
This project has some examples of Go Codes to use as AWS Lambda functions

## 1. Architecture
This example uses the following technologies:
- Go

Other Keywords:
- AWS RDS;
- AWS DynamoDB;
- AWS SES;
- AWS Api Gateway;

## 2. Generating de zip to upload into AWS Lambda:
1. Please follow the steps provided by Go Official page to install Go: https://golang.org/doc/install?download=go1.12.7.linux-amd64.tar.gz
2. For every '.go' file inside lambdas folder run the 'get command' to get the dependencies. Example.: For 'apiGatewayLambda.go' run the following commands:
    
    ```
    go get github.com/gorilla/mux
    go get go.mongodb.org/mongo-driver
    ```

3. Inside project folder run the command to run the project:
    
    ```
    GOOS=linux GOARCH=amd64 go build -o main apiGatewayLambda.go
    zip deployment.zip main
    rm main
    ```
    
    **Obs.:** The deployment commands can be executed throw the 'deploymentScript.sh' file too.

4. Upload de zip file at AWS Lambda;

## 3. References
This project was build following the instructions and documentations provided by this pages:
- https://golang.org/ (Last visited in: 29/07/2019)
- https://medium.com/@rafaelacioly/construindo-uma-api-restful-com-go-d6007e4faff6 (Last visited in: 29/07/2019)
- https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b (Last visited in: 29/07/2019)
- https://www.thepolyglotdeveloper.com/2019/02/developing-restful-api-golang-mongodb-nosql-database/ (Last visited in: 29/07/2019)
- https://www.javacodegeeks.com/2018/11/build-restful-api-go-using-aws-lambda.html (Last visited in: 29/07/2019)