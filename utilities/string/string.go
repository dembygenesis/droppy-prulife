package string

import (
	// "fmt"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"os"
	"strings"
	"time"
)

/**
	This utility handles the authentication via JWT.
	Base64
 */

const (
	MinCost     int = 2  // the minimum allowable cost as passed in to GenerateFromPassword
	MaxCost     int = 31 // the maximum allowable cost as passed in to GenerateFromPassword
	DefaultCost int = 10 // the cost that will actually be set if a cost below MinCost is passed into GenerateFromPassword

)

var timeFormat = "2006-01-02 15:04:05"

// Encrypts as password using bcrypt
func Encrypt(text string) (string, error) {
	var encryptedPassword = ""

	data, err := bcrypt.GenerateFromPassword([]byte(text), 10)

	if err != nil {
		return encryptedPassword, err
	}

	encryptedPassword = string(data)

	fmt.Println("encryptedPassword", encryptedPassword)

	return encryptedPassword, err
}

// Attempts to match a bcrypt encrypted password
func Decrypt(dataEncrypted string, regularData string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dataEncrypted), []byte(regularData))

	if err == nil {
		return true
	} else {
		return false
	}
}

func MakeJWT(i interface{}) (string, error) {
	userId := fmt.Sprintf("%v", i)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := []byte(os.Getenv("SECRET_KEY"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"timestamp": time.Now().Format(timeFormat),
	})

	tokenString, err := token.SignedString(secretKey)

	return tokenString, nil
}

func ParseJWT(i interface{})  (jwt.MapClaims, error) {
	parsedString := fmt.Sprintf("%v", i)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := []byte(os.Getenv("SECRET_KEY"))

	token, err := jwt.Parse(parsedString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
		// fmt.Println("order of the quath", claims["userId"], claims["timestamp"])
	} else {
		return nil, err
	}
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func GetJoinedStringForWhereIn(s []string) string {
	var joinedStr string

	for i, _ := range s {
		joinedStr = joinedStr + "\"" + s[i] + "\","
		// joinedStr = joinedStr + "" + s[i] + ","
	}

	joinedStr = TrimSuffix(joinedStr, ",")

	return joinedStr
}