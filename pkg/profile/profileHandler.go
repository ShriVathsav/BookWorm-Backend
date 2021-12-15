package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"the-book-store/db"
	"the-book-store/helpers"
	"the-book-store/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorFailedToFetchRecord     = "failed to fetch record"
	ErrorInvalidData             = "invalid data"
	ErrorCouldNotMarshalItem     = "could not marshal item"
	ErrorCouldNotDeleteItem      = "could not delete item"
	ErrorCouldNotUpdateItem      = "could not update item"
)

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

// GET profile/{profileId}
func GetProfileHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	profileId := req.PathParameters["profileId"]
	var profile models.Profile
	err := GetProfile(profileId, &profile)

	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}

	return helpers.ApiResponse(http.StatusOK, profile)
}

// GET profile/byCognitoId/{cognitoId}
func GetProfileByCognitoIdHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	cognitoId := req.PathParameters["cognitoId"]
	var profile models.Profile
	err := GetProfileByCognitoId(cognitoId, &profile)

	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}

	return helpers.ApiResponse(http.StatusOK, profile)
}

// POST profile/
func CreateProfileHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	var profile models.Profile
	if err := json.Unmarshal([]byte(req.Body), &profile); err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	// fmt.Println(profile, r.Body)
	err := CreateProfile(profile)
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusCreated, profile)
}

// PUT profile/{profileId}
func UpdateProfileHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	var profile models.Profile
	if err := json.Unmarshal([]byte(req.Body), &profile); err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	profileId := req.PathParameters["profileId"]
	err := UpdateProfile(profileId, profile)
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusOK, profileId)
}

// PUT profile/{profileId}/updateCart
func UpdateCartHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	var cart []models.CartItem
	profileId := req.PathParameters["profileId"]
	if err := json.Unmarshal([]byte(req.Body), &cart); err != nil {
		return nil, errors.New(ErrorInvalidData)
	}
	err := UpdateCart(profileId, cart)
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusOK, profileId)
}

// PUT profile/{profileId}/updateProfileImage
func UpdateProfileImageHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	var profile models.Profile
	if err := json.Unmarshal([]byte(req.Body), &profile); err != nil {
		return nil, errors.New(ErrorInvalidData)
	}
	profileId := req.PathParameters["profileId"]
	err := UpdateProfileImage(profileId, profile)
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusOK, profileId)
}

// DELETE profile/{profileId}
func DeleteProfileHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	profileId := req.PathParameters["profileId"]
	err := DeleteProfile(profileId)
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusOK, profileId)
	// json.NewEncoder(w).Encode("profile not found")

}

// get all profiles from the DB and return it
func GetAllProfiles() ([]primitive.M, error) {
	cur, err := db.DatabaseObj.Collection("profile").Find(context.Background(), bson.D{{}})
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		// fmt.Println("cur..>", cur, "result", reflect.TypeOf(result), reflect.TypeOf(result["_id"]))
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results, nil
}

func GetProfile(profileId string, profile *models.Profile) error {
	id, _ := primitive.ObjectIDFromHex(profileId)
	fmt.Println("Object id", id, profileId)

	//collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := db.DatabaseObj.Collection("profile").FindOne(ctx, bson.M{"_id": id}).Decode(profile)
	if err != nil {
		return errors.New(ErrorFailedToFetchRecord)
	}
	fmt.Println("profile", profile)
	return err
}

func GetProfileByCognitoId(cognitoId string, profile *models.Profile) error {
	fmt.Println("Cognito id", cognitoId)

	//collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := db.DatabaseObj.Collection("profile").FindOne(ctx, bson.M{"cognitoid": cognitoId}).Decode(profile)
	fmt.Println("profile", profile)
	return err
}

// Insert one profile in the DB
func CreateProfile(profile models.Profile) error {
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()
	insertResult, err := db.DatabaseObj.Collection("profile").InsertOne(context.Background(), profile)

	fmt.Println(profile)
	//mongoId.(primitive.ObjectID).Hex()
	if err != nil {
		return errors.New(ErrorCouldNotUpdateItem)
	}

	profile.ID, _ = primitive.ObjectIDFromHex(insertResult.InsertedID.(primitive.ObjectID).Hex())
	fmt.Println("Inserted a Single Record ", insertResult.InsertedID, *insertResult)
	return nil
}

// profile Update method, update profile
func UpdateProfile(profileId string, profile models.Profile) error {
	fmt.Println(profileId)
	id, _ := primitive.ObjectIDFromHex(profileId)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"phone":         profile.Phone,
		"address1":      profile.Address1,
		"address2":      profile.Address2,
		"profile_image": profile.ProfileImage,
		"pincode":       profile.Pincode,
		"updated_at":    time.Now(),
	},
	}
	result, err := db.DatabaseObj.Collection("profile").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return errors.New(ErrorCouldNotUpdateItem)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
	return nil
}

// profile Image Update method, update profile
func UpdateProfileImage(profileId string, profile models.Profile) error {
	fmt.Println(profileId)
	id, _ := primitive.ObjectIDFromHex(profileId)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"profile_image": profile.ProfileImage,
	},
	}
	result, err := db.DatabaseObj.Collection("profile").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return errors.New(ErrorCouldNotUpdateItem)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
	return nil
}

// profile undo method, update profile's status to false
func UpdateCart(profileId string, cart []models.CartItem) error {
	fmt.Println(profileId)
	id, _ := primitive.ObjectIDFromHex(profileId)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"cart": cart}}
	result, err := db.DatabaseObj.Collection("profile").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return errors.New(ErrorCouldNotUpdateItem)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
	return nil
}

// delete one profile from the DB, delete by ID
func DeleteProfile(profile string) error {
	fmt.Println(profile)

	id, _ := primitive.ObjectIDFromHex(profile)
	filter := bson.M{"_id": id}

	d, err := db.DatabaseObj.Collection("profile").DeleteOne(context.Background(), filter)
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}

	fmt.Println("Deleted Document", d.DeletedCount)
	return nil
}
