package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"the-book-store/db"
	"the-book-store/dtos"
	"the-book-store/helpers"
	"the-book-store/models"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"

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

// GET order/getAllByProfile/{profileId}
func GetAllOrdersHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {
	statusValuesRaw := req.QueryStringParameters["statusValues"]
	statusValues := strings.Split(statusValuesRaw, ",")
	fmt.Println(statusValues, req.MultiValueQueryStringParameters)
	profileId := req.PathParameters["profileId"]
	payload, err := GetAllOrders(profileId, statusValues)
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusOK, payload)
}

// GET order/getAllWaiting/{profileId}
func GetAllWaitingOrdersHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {
	statusValuesRaw := req.QueryStringParameters["statusValues"]
	statusValues := strings.Split(statusValuesRaw, ",")
	fmt.Println(statusValues, "HELLO I`M STATUS VALUES")
	fmt.Println(statusValues, req.MultiValueQueryStringParameters["statusValues"])
	profileId := req.PathParameters["profileId"]
	payload, err := GetAllWaitingOrders(profileId, statusValues)
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusOK, payload)
}

// GET order/{orderId}
func GetOrderHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	orderId := req.PathParameters["orderId"]
	var order models.Order
	err := GetOrder(orderId, &order)

	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}

	return helpers.ApiResponse(http.StatusOK, order)
}

// POST order/
func CreateOrderHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	var orders []models.Order
	var payment dtos.Payment
	if err := json.Unmarshal([]byte(req.Body), &payment); err != nil {
		fmt.Println(err, "PAYMENT ERROR 1")
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	fmt.Println(payment, req.Body)
	paymentError := Payment(&payment)
	if paymentError != nil {
		fmt.Println(paymentError, "PAYMENT ERROR 2")
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(paymentError.Error()),
		})
	}
	orderError := CreateOrder(payment.Orders)
	if orderError != nil {
		fmt.Println(orderError, "PAYMENT ERROR 3")
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(orderError.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusCreated, orders)
}

// PUT order/{orderId}/updateStatus
func UpdateOrderStatusHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	var order models.Order
	if err := json.Unmarshal([]byte(req.Body), &order); err != nil {
		fmt.Println(err, "UPDATE ORDER 1 ERROR")
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	orderIdRaw := req.PathParameters["orderId"]
	err := UpdateOrderStatus(orderIdRaw, order)
	if err != nil {
		fmt.Println(err, "UPDATE ORDER 2 ERROR")
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusOK, orderIdRaw)
}

// DELETE order/{orderId}
func DeleteOrderHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	orderId := req.PathParameters["orderId"]
	err := DeleteOrder(orderId)
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusOK, orderId)
	// json.NewEncoder(w).Encode("order not found")

}

// get all profiles from the DB and return it
func GetAllOrders(profileId string, statusValues []string) ([]primitive.M, error) {
	fmt.Println(profileId, "BUYER ORDER")
	fmt.Println(statusValues, "STATUS VALUES PLAIN")
	cur, err := db.DatabaseObj.Collection("order").Find(context.Background(), bson.M{
		"buyer":  profileId,
		"status": bson.M{"$in": statusValues},
	})
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

		var book models.Book
		GetBook(result["book"].(string), &book)
		result["book"] = book
		// fmt.Println("cur..>", cur, "result", reflect.TypeOf(result), reflect.TypeOf(result["_id"]))
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results, nil
}

func GetBook(bookId string, book *models.Book) error {
	id, _ := primitive.ObjectIDFromHex(bookId)
	fmt.Println("Object id", id, bookId)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := db.DatabaseObj.Collection("book").FindOne(ctx, bson.M{"_id": id}).Decode(book)
	if err != nil {
		return errors.New(ErrorFailedToFetchRecord)
	}

	fmt.Println("book", book)
	return err
}

// get all profiles from the DB and return it
func GetAllWaitingOrders(profileId string, statusValues []string) ([]primitive.M, error) {
	fmt.Println(statusValues, "STATUS VALUES WAITING")
	cur, err := db.DatabaseObj.Collection("order").Find(context.Background(), bson.M{
		"seller": profileId,
		"status": bson.M{"$in": statusValues},
		/*
			"$or": []interface{}{
				bson.M{"status": bson.A{"DELIVERED", "IN PROGRESS"}},
				bson.M{"status": "COLLECTED"},
			},
		*/
	})
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
		var book models.Book
		GetBook(result["book"].(string), &book)
		result["book"] = book
		//fmt.Println(result, "order waiting")
		// fmt.Println("cur..>", cur, "result", reflect.TypeOf(result), reflect.TypeOf(result["_id"]))
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results, nil
}

func GetOrder(orderId string, order *models.Order) error {
	id, _ := primitive.ObjectIDFromHex(orderId)
	fmt.Println("Object id", id, orderId)

	//collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := db.DatabaseObj.Collection("order").FindOne(ctx, bson.M{"_id": id}).Decode(order)
	fmt.Println("order", order)
	if err != nil {
		return errors.New(ErrorFailedToFetchRecord)
	}
	return nil
}

// Insert one order in the DB
func CreateOrder(orders []models.Order) error {
	for index, order := range orders {
		fmt.Println("At index --- ", index, "order value is --- ", order)
		order.CreatedAt = time.Now()
		order.UpdatedAt = time.Now()
		insertResult, err := db.DatabaseObj.Collection("order").InsertOne(context.Background(), order)
		if err != nil {
			return errors.New(ErrorCouldNotUpdateItem)
		}
		fmt.Println("Inserted a Single Record ", insertResult.InsertedID, *insertResult)

		order.ID, _ = primitive.ObjectIDFromHex(insertResult.InsertedID.(primitive.ObjectID).Hex())

		var book models.Book
		UpdateBookQuantityAfterOrder(order.Book, book, order.Quantity)
	}
	return nil
}

func UpdateBookQuantityAfterOrder(bookId string, book models.Book, orderedQuantity int64) error {
	fmt.Println(bookId)
	id, _ := primitive.ObjectIDFromHex(bookId)

	// GET BOOK BY ID
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := db.DatabaseObj.Collection("book").FindOne(ctx, bson.M{"_id": id}).Decode(&book)
	if err != nil {
		return errors.New(ErrorFailedToFetchRecord)
	}

	bookQuantity := book.StocksLeft
	filter := bson.M{"_id": id}
	fmt.Println(bookQuantity-orderedQuantity > 0, bookQuantity, orderedQuantity, "PRINTING QUANTITRIES")
	finalQuantity := bookQuantity - orderedQuantity
	update := bson.M{"$set": bson.M{"stocksleft": finalQuantity, "instock": bookQuantity-orderedQuantity > 0}}
	result, err := db.DatabaseObj.Collection("book").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return errors.New(ErrorCouldNotUpdateItem)
	}

	fmt.Println("modified count: ", result, result.ModifiedCount)
	return nil
}

// Order Update method, update order
func UpdateOrderStatus(orderId string, order models.Order) error {
	fmt.Println(orderId)
	id, _ := primitive.ObjectIDFromHex(orderId)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"status":       order.Status,
		"deliverydate": order.DeliveryDate,
		"updated_at":   time.Now(),
	},
	}
	result, err := db.DatabaseObj.Collection("order").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return errors.New(ErrorCouldNotUpdateItem)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
	return nil
}

// delete one profile from the DB, delete by ID
func DeleteOrder(order string) error {
	fmt.Println(order)

	id, _ := primitive.ObjectIDFromHex(order)
	filter := bson.M{"_id": id}

	d, err := db.DatabaseObj.Collection("order").DeleteOne(context.Background(), filter)
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}

	fmt.Println("Deleted Document", d.DeletedCount)
	return nil
}

func Payment(payment *dtos.Payment) error {

	//apiKey := "sk_test_51HHWVYBB3FZLOZ1moSSahrY2wOpWXWmAsDNWHBJzxUJSocWEACxNs3e2SwIXVFJlV3Sp1HsWcinVpuA4xs0X5kMg00VQ29dSwD"
	apiKey := "sk_test_dp3Wxv3ZhwxyQvwxks34udKh005Ueb8ZLy"
	fmt.Println(apiKey + "asdasd")
	stripe.Key = apiKey
	_, err := charge.New(&stripe.ChargeParams{
		Amount:       stripe.Int64(payment.TotalAmount),
		Currency:     stripe.String(string(stripe.CurrencyINR)),
		Description:  stripe.String(payment.Description),
		Source:       &stripe.SourceParams{Token: stripe.String(payment.StripeToken)},
		ReceiptEmail: stripe.String(payment.ReceiptEmail)})

	return err
}
