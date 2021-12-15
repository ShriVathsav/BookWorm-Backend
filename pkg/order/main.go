package main

import (
	"fmt"
	"the-book-store/db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	fmt.Println("Entering MAIN")
	//region := os.Getenv("AWS_REGION")
	fmt.Println("BEFORE BOOK HANDLER")
	lambda.Start(handler)
	fmt.Println("Exiting MAIN")
}

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return MatchRouteOrder(req)
}

func init() {
	fmt.Println("INITIALIZING DATABASE")
	db.Init()
	fmt.Println("INITIALIZED DATABASE")
}
