# GO REST API Example
This is a example of REST API using GO

## 1. Architecture
This example uses the following technologies:
- Go
- Postgresql

Frameworks used:
- gorilla/mux
- go.mongodb.org/mongo-driver

## 2. Installing amd running
1. Please follow the steps provided by Go Official page to install Go: https://golang.org/doc/install?download=go1.12.7.linux-amd64.tar.gz
2. Inside project folder run the command: 
    ```
    go get github.com/gorilla/mux
    go get go.mongodb.org/mongo-driver
    ```
3. Inside project folder run the command to run the project:
    ```
    go build && ./go-rest-api-case
    ```

## 3. References
This project was build following the instructions and documentations provided by this pages:
- https://golang.org/ (Last visited in: 29/07/2019)
- https://medium.com/@rafaelacioly/construindo-uma-api-restful-com-go-d6007e4faff6 (Last visited in: 29/07/2019)
- https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b (Last visited in: 29/07/2019)
- https://www.thepolyglotdeveloper.com/2019/02/developing-restful-api-golang-mongodb-nosql-database/ (Last visited in: 29/07/2019)