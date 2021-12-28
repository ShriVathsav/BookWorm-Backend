package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ToDoList struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Task   string             `json:"task,omitempty"`
	Status bool               `json:"status,omitempty"`
}

type Book struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title           string             `json:"title,omitempty"`
	Price           float64            `json:"price,omitempty"`
	SellingPrice    float64            `bson:"selling_price" json:"selling_price,omitempty"`
	Category        string             `json:"category,omitempty"`
	Description     string             `json:"description,omitempty"`
	Dimensions      int64              `json:"dimensions,omitempty"`
	NumberOfPages   int64              `bson:"number_of_pages" json:"number_of_pages,omitempty"`
	BookType        string             `bson:"book_type" json:"book_type,omitempty"`
	Author          string             `json:"author,omitempty"`
	Year            int16              `json:"year,omitempty"`
	Weight          float64            `json:"weight,omitempty"`
	Condition       string             `json:"condition,omitempty"`
	Publisher       string             `json:"publisher,omitempty"`
	StocksLeft      int64              `bson:"stocks_left" json:"stocks_left,omitempty"`
	DeliveryTime    int64              `bson:"delivery_time" json:"delivery_time,omitempty"`
	CountryOfOrigin string             `bson:"country_of_origin" json:"country_of_origin,omitempty"`
	Language        string             `json:"language,omitempty"`
	Status          string             `json:"status,omitempty"`
	AverageRating   float64            `bson:"average_rating" json:"average_rating"`
	ReviewCount     int64              `bson:"review_count" json:"review_count"`
	FiveStar        int64              `bson:"five_star" json:"five_star"`
	FourStar        int64              `bson:"four_star" json:"four_star"`
	ThreeStar       int64              `bson:"three_star" json:"three_star"`
	TwoStar         int64              `bson:"two_star" json:"two_star"`
	OneStar         int64              `bson:"one_star" json:"one_star"`
	InStock         bool               `bson:"in_stock" json:"in_stock,omitempty"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
	Profile         string             `json:"profile,omitempty"`
	Reviews         []string           `json:"review,omitempty"`
	Images          []string           `json:"images,omitempty"`
	CoverImage      string             `json:"coverimage,omitempty"`
	PeopleBought    []string           `bson:"people_bought" json:"people_bought,omitempty"`
}

type Profile struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CognitoId     string             `bson:"cognito_id" json:"cognito_id,omitempty"`
	Username      string             `json:"username,omitempty"`
	Email         string             `json:"email,omitempty"`
	ProfileImage  string             `bson:"profile_image" json:"profile_image,omitempty"`
	Phone         string             `json:"phone,omitempty"`
	Address1      string             `json:"address1,omitempty"`
	Address2      string             `json:"address2,omitempty"`
	Pincode       string             `json:"pincode,omitempty"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
	PostedBooks   []string           `json:"posted_books" json:"posted_books,omitempty"`
	Orders        []string           `json:"orders,omitempty"`
	OrdersWaiting []string           `bson:"orders_waiting" json:"orders_waiting,omitempty"`
	Cart          []CartItem         `json:"cart,omitempty"`
}

type Order struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	DeliveryDate string             `bson:"delivery_date" json:"delivery_date,omitempty"`
	Seller       string             `json:"seller,omitempty"`
	Buyer        string             `json:"buyer,omitempty"`
	Book         string             `json:"book,omitempty"`
	Quantity     int64              `json:"quantity,omitempty"`
	Amount       float64            `json:"amount,omitempty"`
	Status       string             `json:"status,omitempty"`
	Reviewed     bool               `json:"reviewed,omitempty"`
	BuyerName    string             `bson:"buyer_name" json:"buyer_name,omitempty"`
	BuyerEmail   string             `bson:"buyer_email" json:"buyer_email,omitempty"`
	Phone        string             `json:"phone,omitempty"`
	Address1     string             `json:"address1,omitempty"`
	Address2     string             `json:"address2,omitempty"`
	Pincode      string             `json:"pincode,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
}

type Review struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Content   string             `json:"content,omitempty"`
	Stars     int32              `json:"stars,omitempty"`
	Images    []string           `json:"images,omitempty"`
	Profile   string             `json:"profile,omitempty"`
	Book      string             `json:"book,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
}

type Cart struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Profile   *Profile           `json:"profile,omitempty"`
	CartItems []*CartItem        `json:"cart_items" json:"cart_items,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
}

type CartItem struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Book      string             `json:"book,omitempty"`
	Quantity  int64              `json:"quantity,omitempty"`
	Amount    float64            `json:"amount,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
}
