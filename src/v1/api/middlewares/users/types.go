package users

/**
Request parameter structs
*/

type ParamsLogin struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}