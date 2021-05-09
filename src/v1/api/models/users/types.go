package users

import "github.com/dembygenesis/droppy-prulife/utilities/response_builder"

type User struct {
	ID        int    `json:"id,omitempty" db:"id"`
	Email     string `json:"email,omitempty" db:"email"`
	Password  string `json:"password,omitempty" db:"password"`
	FirstName string `json:"firstname,omitempty" db:"firstname"`
	LastName  string `json:"lastname,omitempty" db:"lastname"`
	Role      string `json:"role,omitempty" db:"role"`

	Token string `json:"token"`
	/*

		MobileNumber string `json:"mobile_number,omitempty" db:"mobile_number"`

		UserTypeId string `json:"user_type_id,omitempty" db:"user_type_id"`
		CreatedBy string `json:"created_by,omitempty" db:"created_by"`
		CreatedDate string `json:"created_date,omitempty" db:"created_date"`
		LastUpdated string `json:"last_updated,omitempty" db:"last_updated"`
		UpdatedBy string `json:"updated_by,omitempty" db:"updated_by"`
		IsActive string `json:"is_active,omitempty" db:"is_active"`
		BankTypeId string `json:"bank_type_id,omitempty" db:"bank_type_id"`
		BankNo string `json:"bank_no,omitempty" db:"bank_no"`
		Address string `json:"address,omitempty" db:"address"`
		Birthday string `json:"birthday,omitempty" db:"birthday"`
		Gender string `json:"gender,omitempty" db:"gender"`
		M88Account string `json:"m88_account,omitempty" db:"m88_account"`*/
}

type UserListDisplay struct {
	ID           uint   `json:"id" db:"id"`
	FirstName    string `json:"firstname" db:"firstname"`
	LastName     string `json:"lastname" db:"lastname"`
	Email        string `json:"email" db:"email"`
	MobileNumber string `json:"mobile_number" db:"mobile_number"`
	Role         string `json:"role" db:"role"`
	BankType     string `json:"bank_type" db:"bank_type"`
	BankNo       string `json:"bank_no" db:"bank_no"`
	Address      string `json:"address" db:"address"`
	Birthday     string `json:"birthday" db:"birthday"`
	Gender       string `json:"gender" db:"gender"`
	M88Account   string `json:"m88_account" db:"m88_account"`
}

type UserLogin struct {
	ID       uint   `json:"id,omitempty" db:"id"`
	Email    string `json:"email,omitempty" db:"email"`
	Password string `json:"password,omitempty" db:"password"`
}

/**
Parameter Types
*/

type ParamsInsert struct {
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	Email        string `json:"email"`
	MobileNumber string `json:"mobile_number"`
	Password     string `json:"password"`
	UserTypeId   int    `json:"user_type_id"`
	CreatedBy    int
	BankTypeId   int    `json:"bank_type_id"`
	BankNo       string `json:"bank_no"`
	Address      string `json:"address"`
	Birthday     string `json:"birthday"`
	M88Account   string `json:"m88_account"`
	Gender       string `json:"gender"`
	RegionId     int    `json:"region_id"`
}

type ParamsUpdate struct {
	ID           int    `json:"id"`
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	Email        string `json:"email"`
	MobileNumber string `json:"mobile_number"`
	Password     string `json:"password"`
	UserTypeId   int    `json:"user_type_id"`
	BankTypeId   int    `json:"bank_type_id"`
	BankNo       string `json:"bank_no"`
	Address      string `json:"address"`
	Birthday     string `json:"birthday"`
	M88Account   string `json:"m88_account"`
	Gender       string `json:"gender"`
}

type ParamsLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ParamsDelete struct {
	ID int `json:"id"`
}

/**
Response Types
*/

type ResponseUserList struct {
	response_builder.Response
	Data *[]UserListDisplay `json:"data,omitempty"`
}

type ResponseUserSingleDisplay struct {
	ID           string `json:"id" db:"id"`
	FirstName    string `json:"firstname" db:"firstname"`
	LastName     string `json:"lastname" db:"lastname"`
	Email        string `json:"email"  db:"email"`
	MobileNumber string `json:"mobile_number" db:"mobile_number"`
	Role         string `json:"role" db:"role"`
	BankType     string `json:"bank_type" db:"bank_type"`
	BankNo       string `json:"bank_no" db:"bank_no"`
	Address      string `json:"address" db:"address"`
	Birthday     string `json:"birthday" db:"birthday"`
	Gender       string `json:"gender" db:"gender"`
	UserTypeId   int    `json:"user_type_id" db:"user_type_id"`
	BankTypeId   int    `json:"bank_type_id" db:"bank_type_id"`
	M88Account   string `json:"m88_account" db:"m88_account"`
	RegionId     int    `json:"region_id" db:"region_id"`
}

type ResponseLogin struct {
	response_builder.Response
}

type ResponseUserDetailsDisplay struct {
	response_builder.Response
	Data struct {
		UserInfo struct {
			Token       string                `json:"token,omitempty"`
			UserDetails ResponseLoginUserInfo `json:"userDetails,omitempty"`
		} `json:"userInfo,omitempty"`
	} `json:"data,omitempty"`
}

type ResponseLoginUserInfo struct {
	ID        int     `json:"id" db:"id"`
	FirstName string  `json:"firstname" db:"firstname"`
	LastName  string  `json:"lastname" db:"lastname"`
	Role      string  `json:"role" db:"role"`
	Balance   float64 `json:"balance" db:"balance"`
	Region    string  `json:"region" db:"region"`
}

type UserMiddlewareDetails struct {
	UserCount int    `json:"user_count,omitempty" db:"user_count"`
	UserType  string `json:"user_type,omitempty" db:"user_type"`
}
