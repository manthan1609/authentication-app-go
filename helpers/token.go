package helpers

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/manthan1609/authentication-app/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GenerateAllTokens(email string, firstName string, lastName string, userType string, uid string) (signedToken string, signedRefreshToken string, err error) {

	// var userCollection *mongo.Collection = database.OpenCollection("user")

	var SECRET_KEY string = os.Getenv("SECRET_KEY")

	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		UserType:  userType,
		Uid:       uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(24))),
		},
	}

	refreshClaims := &SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(365))),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func ValidateToken(signedToken string) (claims *SignedDetails, err error) {
	var SECRET_KEY string = os.Getenv("SECRET_KEY")

	token, errr := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if errr != nil {
		err = errr
		return
	}

	claims, ok := token.Claims.(*SignedDetails)

	if !ok {
		err = errors.New("token is invalid")
		return
	}

	expired, errr := claims.GetExpirationTime()

	if errr != nil {
		err = errr
		return
	}

	if expired.Unix() < time.Now().Unix() {
		err = errors.New("token is expired")
		return
	}

	return claims, err

}

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	UserType  string
	jwt.RegisteredClaims
}

func UpdateAllTokens(accessToken string, refreshToken string, userId string) bool {
	var userCollection *mongo.Collection = database.OpenCollection("user")

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: accessToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: refreshToken})

	updatedAt, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	if err != nil {
		return false
	}

	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: updatedAt})

	upsert := true
	filter := bson.M{
		"user_id": userId,
	}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err = userCollection.UpdateOne(ctx, filter, updateObj, &opt)

	return err == nil
}
