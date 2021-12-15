package main

import (
	"fmt"
	"the-book-store/helpers"

	"github.com/aws/aws-lambda-go/events"
)

func MatchRouteProfile(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	fmt.Println("hello I`m inside the PROFILE handler")
	fmt.Println(req)
	fmt.Printf("%+v\n", req)
	switch req.HTTPMethod {
	case "GET":
		if req.Resource == "/profile/{profileId}" {
			return GetProfileHandler(req)
		} else if req.Resource == "/profile/getByCognitoId/{cognitoId}" {
			return GetProfileByCognitoIdHandler(req)
		} else {
			return helpers.UnhandledMethod()
		}
	case "POST":
		return CreateProfileHandler(req)
	case "PUT":
		return UpdateProfileHandler(req)
	case "PUTS":
		if req.Resource == "/profile/{profileId}" {
			return UpdateProfileHandler(req)
		} else if req.Resource == "/profile/{profileId}/updateCart" {
			return UpdateCartHandler(req)
		} else if req.Resource == "/profile/{profileId}/updateProfileImage" {
			return UpdateProfileImageHandler(req)
		} else {
			return helpers.UnhandledMethod()
		}
	case "DELETE":
		return DeleteProfileHandler(req)
	default:
		fmt.Println("Exiting handler")
		return helpers.UnhandledMethod()
	}
}
