package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://bookappshriv:bookappshriv@cluster0.jaqan.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

// Database Name
const DbName = "bookWormDB"

// collection object/instance
var DatabaseObj *mongo.Database

// create connection with mongo db
func Init() {

	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	DatabaseObj = client.Database(DbName)

	fmt.Printf("DatabaseObj = %T\n", DatabaseObj)

	fmt.Println("Collection instance created!")
}

/*
func GetAllReviews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//params := mux.Vars(r)
	//payload := getAllReviews(params["id"])
	//json.NewEncoder(w).Encode(payload)
}
*/
