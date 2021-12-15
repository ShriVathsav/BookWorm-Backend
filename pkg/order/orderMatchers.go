package main

import (
	"fmt"
	"the-book-store/helpers"

	"github.com/aws/aws-lambda-go/events"
)

func MatchRouteOrder(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	fmt.Println("hello I`m inside the ORDER handler")
	fmt.Println(req)
	fmt.Printf("%+v\n", req)
	switch req.HTTPMethod {
	case "GET":
		if req.Resource == "/order/getAllByProfile/{profileId}" {
			return GetAllOrdersHandler(req)
		} else if req.Resource == "/order/getAllWaiting/{profileId}" {
			return GetAllWaitingOrdersHandler(req)
		} else if req.Resource == "/order/{orderId}" {
			return GetOrderHandler(req)
		} else {
			return helpers.UnhandledMethod()
		}
	case "POST":
		return CreateOrderHandler(req)
	case "PUT":
		return UpdateOrderStatusHandler(req)
	case "DELETE":
		return DeleteOrderHandler(req)
	default:
		fmt.Println("Exiting handler")
		return helpers.UnhandledMethod()
	}
}
