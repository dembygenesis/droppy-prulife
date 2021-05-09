package users

import (
	// "errors"
	"errors"
	"fmt"

	UserModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/users"

	// "fmt"

	StringUtility "github.com/dembygenesis/droppy-prulife/utilities/string"
)

func GetAll() (*[]UserModel.UserListDisplay, error) {
	user := UserModel.User{}

	res, err := user.GetAll()

	return res, err
}

// Returns a token if successful, error if not lol
func Login(email string, password string) (string, UserModel.ResponseLoginUserInfo, error) {

	var responseLoginUserInfo UserModel.ResponseLoginUserInfo

	var jwtToken = ""

	user := UserModel.User{Email: email, Password: password}

	// Get password via email
	password, id, err := user.GetPasswordAndIdByEmail()

	if err != nil {
		return jwtToken, responseLoginUserInfo, errors.New("something went wrong when trying to check the password")
	}

	// Attempt to match password
	matched := StringUtility.Decrypt(password, user.Password)

	if matched == false {
		return jwtToken, responseLoginUserInfo, errors.New("username/password match failed")
	}

	// Extract JWT
	jwtToken, err = StringUtility.MakeJWT(id)

	if err != nil {
		return jwtToken, responseLoginUserInfo, errors.New("something went wrong when trying to extract a JWTToken")
	}

	// Get login details
	responseLoginUserInfo, err = user.GetLoginDetails()

	if err != nil {
		fmt.Println("hhohoh", err)
		// return jwtToken, responseLoginUserInfo, errors.New("something went wrong when trying to get report login user info")
		fmt.Println("password", password)
		return jwtToken, responseLoginUserInfo, err
	}

	fmt.Print(user)

	return jwtToken, responseLoginUserInfo, nil
}
