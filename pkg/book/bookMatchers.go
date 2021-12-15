package main

import (
	"fmt"
	"the-book-store/helpers"

	"github.com/aws/aws-lambda-go/events"
)

func MatchRouteBook(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	fmt.Println("hello I`m inside the BOOK handler")
	fmt.Println(req)
	fmt.Printf("%+v\n", req)
	switch req.HTTPMethod {
	case "GET":
		if req.Resource == "/book" {
			return GetAllBooksHandler(req)
		} else if req.Resource == "/book/byProfile/{profileId}" {
			return GetBooksPostedHandler(req)
		} else if req.Resource == "/book/getAllById" {
			return GetBooksByIdHandler(req)
		} else if req.Resource == "/book/{bookId}" {
			return GetBookHandler(req)
		} else if req.Resource == "/book/search" {
			return SearchBooksHandler(req)
		} else {
			return helpers.UnhandledMethod()
		}
	case "POST":
		if req.Resource == "/book" {
			return CreateBookHandler(req)
		} else if req.Resource == "/book/uploadimage" {
			return HandleImageUpload(req)
		} else {
			return helpers.UnhandledMethod()
		}
	case "PUT":
		if req.Resource == "/book/{bookId}" {
			return UpdateBookHandler(req)
		} else if req.Resource == "/book/{bookId}/editStatus" {
			return EditBookStatusHandler(req)
		} else if req.Resource == "/book/{bookId}/editQuantity" {
			return EditBookQuantityHandler(req)
		} else {
			return helpers.UnhandledMethod()
		}
	case "DELETE":
		return DeleteBookHandler(req)
	default:
		fmt.Println("Exiting handler")
		return helpers.UnhandledMethod()
	}
}
