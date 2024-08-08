package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	FirstName    *string            `json:"first_name" validate:"omitempty,required,min=2,max=100" bson:"first_name"`
	LastName     *string            `json:"last_name" validate:"omitempty,required,min=2,max=100" bson:"last_name"`
	Password     *string            `json:"password" validate:"required,min=6" bson:"password"`
	Email        *string            `json:"email" validate:"email,required" bson:"email"`
	Phone        *string            `json:"phone"  bson:"phone" validate:"omitempty,required"`
	Token        *string            `json:"token" bson:"token"`
	UserType     *string            `json:"user_type" validate:"omitempty,required,eq=ADMIN|eq=USER" bson:"user_type"`
	RefreshToken *string            `json:"refresh_token" bson:"refresh_token"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	UserId       string             `json:"user_id" bson:"user_id"`
}
