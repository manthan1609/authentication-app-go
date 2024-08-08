package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")

	err = nil

	if userType == "USER" && uid != userId {
		err = errors.New("unauthorized  access")
		return err
	}

	err = CheckUserType(c, userType)

	return err
}

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")

	err = nil

	if userType != role {
		err = errors.New("unauthorized access")
		return err
	}

	return err
}

func VerifyPassword(hashedPassword string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
