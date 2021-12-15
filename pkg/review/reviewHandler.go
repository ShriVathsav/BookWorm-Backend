package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
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

// GET review/getAllByBook/{bookId}
func GetAllReviewsHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	//params := mux.Vars(r)
	//bookId := r.URL.Query().Get("id")
	bookId := req.PathParameters["bookId"]
	payload, err := getAllReviews(bookId)
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusOK, payload)
}

// GET review/{reviewId}
func GetReviewHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {
	reviewIdRaw := req.PathParameters["reviewId"]
	fmt.Println("Object id", reviewIdRaw)
	var review models.Review
	err := GetReview(reviewIdRaw, &review)
	fmt.Println("Review", review)
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusOK, review)
}

// POST review/
func CreateReviewHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	var review models.Review
	if err := json.Unmarshal([]byte(req.Body), &review); err != nil {
		return nil, errors.New(ErrorInvalidData)
	}
	// fmt.Println(task, r.Body)
	err := insertReview(review)
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusCreated, review)
}

// PUT review/{reviewId}
func UpdateReviewHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	reviewId := req.PathParameters["reviewId"]
	var review models.Review
	if err := json.Unmarshal([]byte(req.Body), &review); err != nil {
		return nil, errors.New(ErrorInvalidData)
	}
	err := updateReview(reviewId, review)
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusOK, reviewId)
}

// DELETE review/{reviewId}
func DeleteReviewHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {
	reviewId := req.PathParameters["reviewId"]
	err := deleteReview(reviewId)
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusOK, reviewId)
	// json.NewEncoder(w).Encode("Task not found")

}

// get all task from the DB and return it
func getAllReviews(bookId string) ([]primitive.M, error) {

	//err := db.DatabaseObj.Collection("book").FindOne(context.Background(), bson.M{{""}})

	//BOOK ID
	//id, _ := primitive.ObjectIDFromHex(bookId)
	fmt.Println(bookId, "BOOK ID FOR REVIEWS")

	cur, err := db.DatabaseObj.Collection("review").Find(context.Background(), bson.M{"book": bookId})

	//GetAllReviewCount(bookId)

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
		var profile models.Profile
		GetProfile(result["profile"].(string), &profile)
		result["profile"] = profile
		// fmt.Println("cur..>", cur, "result", reflect.TypeOf(result), reflect.TypeOf(result["_id"]))
		results = append(results, result)

	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results, nil
}

func GetReview(reviewId string, review *models.Review) error {
	id, _ := primitive.ObjectIDFromHex(reviewId)
	fmt.Println("Object id", id, reviewId)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := db.DatabaseObj.Collection("review").FindOne(ctx, bson.M{"_id": id}).Decode(review)
	if err != nil {
		return errors.New(ErrorFailedToFetchRecord)
	}
	fmt.Println("review", review)
	return err
}

func GetAllReviewCount(bookId string) (int64, error) {
	itemCount, err := db.DatabaseObj.Collection("review").CountDocuments(context.Background(), bson.M{"book": bookId})

	if err != nil {
		return 0, errors.New(ErrorFailedToFetchRecord)
	}
	fmt.Println(itemCount, reflect.TypeOf(itemCount))

	//GetStarReviewCount(bookId)

	return itemCount, nil
}

func GetStarReviewCount(bookId string) (map[string]int64, error) {
	fiveStar, err1 := db.DatabaseObj.Collection("review").CountDocuments(context.Background(), bson.M{"stars": 5, "book": bookId})
	fourStar, err2 := db.DatabaseObj.Collection("review").CountDocuments(context.Background(), bson.M{"stars": 4, "book": bookId})
	threeStar, err3 := db.DatabaseObj.Collection("review").CountDocuments(context.Background(), bson.M{"stars": 3, "book": bookId})
	twoStar, err4 := db.DatabaseObj.Collection("review").CountDocuments(context.Background(), bson.M{"stars": 2, "book": bookId})
	oneStar, err5 := db.DatabaseObj.Collection("review").CountDocuments(context.Background(), bson.M{"stars": 1, "book": bookId})

	fmt.Println(fiveStar, fourStar, threeStar, twoStar, oneStar)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}
	return map[string]int64{"fiveStar": fiveStar, "fourStar": fourStar, "threeStar": threeStar, "twoStar": twoStar, "oneStar": oneStar}, nil
}

// Insert one task in the DB
func insertReview(review models.Review) error {
	insertResult, err := db.DatabaseObj.Collection("review").InsertOne(context.Background(), review)

	if err != nil {
		return errors.New(ErrorCouldNotUpdateItem)
	}

	fmt.Println("Inserted a Single Record ", insertResult.InsertedID)
	return nil

}

func UpdateBookAfterReview(bookId string, review models.Review, updateType string) error {
	fmt.Println(bookId)
	id, _ := primitive.ObjectIDFromHex(bookId)
	var book models.Book

	// GET BOOK BY ID
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := db.DatabaseObj.Collection("book").FindOne(ctx, bson.M{"_id": id}).Decode(&book)

	bookReviewCount := book.ReviewCount
	bookReviewAvg := book.AverageRating * float64(bookReviewCount)
	if updateType == "INSERT" {
		bookReviewCount += 1
	}
	newBookRating := (bookReviewAvg + float64(review.Stars)) / float64(bookReviewCount)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"reviewcount":   bookReviewCount,
		"averagerating": newBookRating,
	}}
	result, err := db.DatabaseObj.Collection("book").UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
	return err
}

// task complete method, update task's status to true
func updateReview(reviewId string, review models.Review) error {
	fmt.Println(reviewId)
	id, _ := primitive.ObjectIDFromHex(reviewId)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"stars": review.Stars, "content": review.Content}}
	result, err := db.DatabaseObj.Collection("review").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return errors.New(ErrorCouldNotUpdateItem)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
	return nil
}

// delete one review from the DB, delete by ID
func deleteReview(reviewId string) error {
	fmt.Println(reviewId)
	id, _ := primitive.ObjectIDFromHex(reviewId)
	filter := bson.M{"_id": id}
	d, err := db.DatabaseObj.Collection("review").DeleteOne(context.Background(), filter)
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}

	fmt.Println("Deleted Document", d.DeletedCount)
	return nil
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
