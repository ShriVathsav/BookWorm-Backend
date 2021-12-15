package dtos

import (
	"the-book-store/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewDto struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Content   string             `json:"content,omitempty"`
	Stars     int32              `json:"stars,omitempty"`
	Images    []string           `json:"images,omitempty"`
	Profile   models.Profile     `json:"profile,omitempty"`
	Book      string             `json:"book,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
}

type Filters struct {
	MinPrice      float64  `json:"minprice,omitempty"`
	MaxPrice      float64  `json:"maxprice,omitempty"`
	Stock         []bool   `json:"stock,omitempty"`
	DeliveryTime  int64    `json:"deliverytime,omitempty"`
	BookCondition []string `json:"bookcondition,omitempty"`
	Rating        float64  `json:"rating,omitempty"`
	BookType      []string `json:"booktype,omitempty"`
}

type Payment struct {
	StripeToken  string         `json:"stripe_token,omitempty"`
	TotalAmount  int64          `json:"total_amount,omitempty"`
	Description  string         `json:"description,omitempty"`
	ReceiptEmail string         `json:"receipt_email,omitempty"`
	Orders       []models.Order `json:"orders,omitempty"`
}
