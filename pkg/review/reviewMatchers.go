package main

import (
	"fmt"
	"the-book-store/helpers"

	"github.com/aws/aws-lambda-go/events"
)

func MatchRouteReview(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	fmt.Println("hello I`m inside the REVIEW handler")
	fmt.Println(req)
	fmt.Printf("%+v\n", req)
	switch req.HTTPMethod {
	case "GET":
		if req.Resource == "/review/getAllByBook/{bookId}" {
			return GetAllReviewsHandler(req)
		} else if req.Resource == "/review/{reviewId}" {
			return GetReviewHandler(req)
		} else {
			return helpers.UnhandledMethod()
		}
	case "POST":
		return CreateReviewHandler(req)
	case "PUT":
		return UpdateReviewHandler(req)
	case "DELETE":
		return DeleteReviewHandler(req)
	default:
		fmt.Println("Exiting handler")
		return helpers.UnhandledMethod()
	}
}
