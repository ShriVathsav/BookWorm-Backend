package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"the-book-store/db"
	"the-book-store/dtos"
	"the-book-store/helpers"
	"the-book-store/models"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	awslambda "github.com/grokify/go-awslambda"

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

// GET book/
func GetAllBooksHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {
	payload, err := GetAllBooks()
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusOK, payload)
}

func apiResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	resp.StatusCode = status

	stringBody, _ := json.Marshal(body)
	resp.Body = string(stringBody)
	return &resp, nil
}

//GET book/search
func SearchBooksHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	searchTerm := req.QueryStringParameters["searchTerm"]
	category := strings.Split(req.QueryStringParameters["category"], ",")
	inStockStrings := strings.Split(req.QueryStringParameters["inStock"], ",")

	var inStock []bool
	for _, j := range inStockStrings {
		conv, _ := strconv.ParseBool(j)
		inStock = append(inStock, conv)
	}
	fmt.Println(inStock, "IN STOCK STATUS VALUES BOOL")

	deliveryTime, _ := strconv.ParseInt(req.QueryStringParameters["deliveryTime"], 10, 64)
	condition := strings.Split(req.QueryStringParameters["condition"], ",")
	rating, _ := strconv.ParseFloat(req.QueryStringParameters["rating"], 64)
	minPrice, _ := strconv.ParseFloat(req.QueryStringParameters["minPrice"], 64)
	maxPrice, _ := strconv.ParseFloat(req.QueryStringParameters["maxPrice"], 64)
	bookType := strings.Split(req.QueryStringParameters["bookType"], ",")

	fmt.Println(deliveryTime, condition, rating, minPrice, maxPrice, bookType, inStock, category, "FILTER PARAMS")
	fmt.Println(reflect.TypeOf(deliveryTime), reflect.TypeOf(condition), reflect.TypeOf(rating),
		reflect.TypeOf(minPrice), reflect.TypeOf(maxPrice), reflect.TypeOf(bookType), reflect.TypeOf(inStock),
		reflect.TypeOf(category), "FILTER PARAMS TYPES")
	filters := dtos.Filters{
		MinPrice:      minPrice,
		MaxPrice:      maxPrice,
		Stock:         inStock,
		DeliveryTime:  deliveryTime,
		BookCondition: condition,
		Rating:        rating,
		BookType:      bookType,
	}

	fmt.Println(searchTerm, "SEARCH TERM")
	payload, err := SearchBooks(searchTerm, category, filters)
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusOK, payload)
}

// GET book/byProfile/{profileId}
func GetBooksPostedHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {
	statusValuesRaw := req.QueryStringParameters["statusValues"]
	statusValues := strings.Split(statusValuesRaw, ",")
	fmt.Println(statusValuesRaw, statusValues, req.QueryStringParameters)
	profileId := req.PathParameters["profileId"]
	payload, err := GetBooksPosted(profileId, statusValues)
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusOK, payload)
}

// GET book/getallById
func GetBooksByIdHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	bookIds := req.QueryStringParameters["bookIds"]
	bookIdsNew := strings.Split(bookIds, ",")
	payload, err := GetBooksById(bookIdsNew)
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusOK, payload)
}

// GET book/{bookId}
func GetBookHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	id := req.PathParameters["bookId"]
	var book models.Book
	err := GetBook(id, &book)

	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}

	return helpers.ApiResponse(http.StatusOK, book)
}

// POST book/
func CreateBookHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	var book models.Book
	//_ = json.NewDecoder(req.Body).Decode(&book)
	if err := json.Unmarshal([]byte(req.Body), &book); err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	// fmt.Println(book, r.Body)
	err := CreateBook(&book)
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusCreated, book)
}

// PUT book/{bookId}
func UpdateBookHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	bookId := req.PathParameters["bookId"]
	var book models.Book
	if err := json.Unmarshal([]byte(req.Body), &book); err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	err := UpdateBook(bookId, book)
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusOK, bookId)
}

// PUT book/{bookId}/editStatus
func EditBookStatusHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	bookId := req.PathParameters["bookId"]
	var book models.Book
	if err := json.Unmarshal([]byte(req.Body), &book); err != nil {
		return nil, errors.New(ErrorInvalidData)
	}
	err := EditBookStatus(bookId, book)
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusOK, bookId)
}

// PUT book/{bookId}/editQuantity
func EditBookQuantityHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	bookId := req.PathParameters["bookId"]
	var book models.Book
	if err := json.Unmarshal([]byte(req.Body), &book); err != nil {
		return nil, errors.New(ErrorInvalidData)
	}
	err := EditBookStatus(bookId, book)
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusOK, bookId)
}

// DELETE book/{bookId}
func DeleteBookHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	//params := mux.Vars(r)
	bookId := req.PathParameters["bookId"]
	DeleteBook(bookId)
	//json.NewEncoder(w).Encode(params["id"])
	return helpers.ApiResponse(http.StatusOK, bookId)
	// json.NewEncoder(w).Encode("Book not found")

}

// DELETE book/
func DeleteAllBooksHandler(req events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse,
	error,
) {

	count, err := DeleteAllBooks()
	if err != nil {
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return helpers.ApiResponse(http.StatusOK, count)
	// json.NewEncoder(w).Encode("Book not found")

}

// get all books from the DB and return it
func GetAllBooks() ([]primitive.M, error) {
	cur, err := db.DatabaseObj.Collection("book").Find(context.Background(), bson.D{{}})
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

func SearchBooks(searchTerm string, categories []string, filters dtos.Filters) ([]primitive.M, error) {
	fmt.Println(filters, searchTerm, "FILTERS STRUCT PRINTING")
	cur, err := db.DatabaseObj.Collection("book").Find(context.Background(),
		bson.M{
			"title":         bson.M{"$regex": ".*" + searchTerm + ".*"},
			"status":        "ACTIVE",
			"category":      bson.M{"$in": categories},
			"instock":       bson.M{"$in": filters.Stock},
			"deliverytime":  bson.M{"$lt": filters.DeliveryTime},
			"condition":     bson.M{"$in": filters.BookCondition},
			"booktype":      bson.M{"$in": filters.BookType},
			"averagerating": bson.M{"$gte": filters.Rating},
			"sellingprice":  bson.M{"$gt": filters.MinPrice, "$lt": filters.MaxPrice},
		},
		//bson.D{{"title", primitive.Regex{Pattern: "bh", Options: ""}}},
	)
	if err != nil {
		log.Fatal(err)
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
	return results, err
}

// get all books from the DB and return it
func GetBooksPosted(profileId string, statusValues []string) ([]primitive.M, error) {
	var statusValuesBool []bool
	for _, j := range statusValues {
		conv, _ := strconv.ParseBool(j)
		statusValuesBool = append(statusValuesBool, conv)
	}
	fmt.Println(statusValuesBool, "STATUS VALUES BOOL")
	cur, err := db.DatabaseObj.Collection("book").Find(context.Background(), bson.M{
		"profile": profileId,
		"instock": bson.M{"$in": statusValuesBool},
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
		// fmt.Println("cur..>", cur, "result", reflect.TypeOf(result), reflect.TypeOf(result["_id"]))
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results, nil
}

// get all books from the DB and return it
func GetBooksById(bookIds []string) ([]primitive.M, error) {

	var bookObjectIds []primitive.ObjectID
	for _, bookId := range bookIds {
		id, _ := primitive.ObjectIDFromHex(bookId)
		bookObjectIds = append(bookObjectIds, id)
	}

	cur, err := db.DatabaseObj.Collection("book").Find(context.Background(), bson.M{"_id": bson.M{"$in": bookObjectIds}})
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

func GetBook(bookId string, book *models.Book) error {
	id, _ := primitive.ObjectIDFromHex(bookId)
	fmt.Println("Object id", id, bookId)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := db.DatabaseObj.Collection("book").FindOne(ctx, bson.M{"_id": id}).Decode(book)
	if err != nil {
		return errors.New(ErrorFailedToFetchRecord)
	}

	reviewCount, err1 := GetAllReviewCount(bookId)
	book.ReviewCount = reviewCount
	reviewStars, err2 := GetStarReviewCount(bookId)
	book.FiveStar = reviewStars["fiveStar"]
	book.FourStar = reviewStars["fourStar"]
	book.ThreeStar = reviewStars["threeStar"]
	book.TwoStar = reviewStars["twoStar"]
	book.OneStar = reviewStars["oneStar"]

	if err1 != nil {
		return errors.New(ErrorFailedToFetchRecord)
	}

	if err2 != nil {
		return errors.New(ErrorFailedToFetchRecord)
	}

	fmt.Println("book", book)
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

// Insert one book in the DB
func CreateBook(book *models.Book) error {
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()
	insertResult, err := db.DatabaseObj.Collection("book").InsertOne(context.Background(), book)

	if err != nil {
		return errors.New(ErrorCouldNotUpdateItem)
	}

	book.ID = insertResult.InsertedID.(primitive.ObjectID)
	fmt.Println(insertResult.InsertedID.(primitive.ObjectID), "hello", insertResult.InsertedID, book.ID)
	fmt.Println("Inserted a Single Record ", reflect.TypeOf(insertResult.InsertedID), insertResult.InsertedID)
	return nil
}

// Book Update method, update book
func UpdateBook(bookId string, book models.Book) error {
	fmt.Println(book)
	id, _ := primitive.ObjectIDFromHex(bookId)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"title":             book.Title,
		"price":             book.Price,
		"selling_price":     book.SellingPrice,
		"category":          book.Category,
		"description":       book.Description,
		"dimensions":        book.Dimensions,
		"number_of_pages":   book.NumberOfPages,
		"book_type":         book.BookType,
		"author":            book.Author,
		"year":              book.Year,
		"weight":            book.Weight,
		"condition":         book.Condition,
		"publisher":         book.Publisher,
		"stocks_left":       book.StocksLeft,
		"delivery_time":     book.DeliveryTime,
		"country_of_origin": book.CountryOfOrigin,
		"language":          book.Language,
		"coverimage":        book.CoverImage,
		"images":            book.Images,
	}}
	result, err := db.DatabaseObj.Collection("book").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return errors.New(ErrorCouldNotUpdateItem)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
	return nil
}

// Book Update Status method, update book
func EditBookStatus(bookId string, book models.Book) error {
	fmt.Println(book)
	id, _ := primitive.ObjectIDFromHex(bookId)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": book.Status}}
	result, err := db.DatabaseObj.Collection("book").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return errors.New(ErrorCouldNotUpdateItem)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
	return nil
}

// Book Update Status method, update book
func EditBookQuantity(bookId string, book models.Book) error {
	fmt.Println(book)
	id, _ := primitive.ObjectIDFromHex(bookId)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"stocksleft":    book.StocksLeft,
		"delivery_time": book.DeliveryTime,
	}}
	result, err := db.DatabaseObj.Collection("book").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return errors.New(ErrorCouldNotUpdateItem)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
	return nil
}

// book undo method, update book's status to false
func undoTask(book string) {
	fmt.Println(book)
	id, _ := primitive.ObjectIDFromHex(book)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"title": "falseShriv"}}
	result, err := db.DatabaseObj.Collection("book").UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
}

// delete one book from the DB, delete by ID
func DeleteBook(book string) error {
	fmt.Println(book)

	id, _ := primitive.ObjectIDFromHex(book)
	filter := bson.M{"_id": id}

	d, err := db.DatabaseObj.Collection("book").DeleteOne(context.Background(), filter)
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}

	fmt.Println("Deleted Document", d.DeletedCount)
	return nil
}

// delete all the books from the DB
func DeleteAllBooks() (int64, error) {

	d, err := db.DatabaseObj.Collection("book").DeleteMany(context.Background(), bson.D{{}}, nil)

	if err != nil {
		return 0, errors.New(ErrorCouldNotDeleteItem)
	}

	fmt.Println("Deleted Document", d.DeletedCount)
	return d.DeletedCount, nil
}

func HandleImageUpload(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	//res := events.APIGatewayProxyResponse{}
	fmt.Println("HELLO I`M INSIDE HANDLE IMAGE UPLOAD FUNCTION")
	r, err := awslambda.NewReaderMultipart(req)
	if err != nil {
		fmt.Println(err, "ERROR PRINTING")
		return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	s, _ := json.MarshalIndent(r, "", "\t")
	fmt.Print(string(s), "PRINTING R NEW READER MULTIPART", r)
	fmt.Printf("%#v\n", r)
	fmt.Println()

	part, err := r.NextPart()
	fmt.Println(part, reflect.TypeOf(part), "PRINTING R.NEXTPART()")
	fmt.Printf("%#v\n", part)
	if err != nil {
		fmt.Println(err, "ERROR PRINTING PART")
	}
	/*
		content, err := ioutil.ReadAll(part)
		if err != nil {
			return helpers.ApiResponse(http.StatusBadRequest, ErrorBody{
				aws.String(err.Error()),
			})
		}
	*/
	return helpers.ApiResponse(http.StatusOK, "IMAGE READ SUCCESSFUL")
}
