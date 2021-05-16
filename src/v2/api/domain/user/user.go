package user

type User struct {
	Id          uint
	Firstname   string
	Lastname    string
	UserType    string //
	UserBalance float64 //
}

type GormComplexQuery struct {
	// 100 fields of different tables
}

// !! DISCONNECTED

// Validates the user email's format is correct
func (u *User) ValidateEmail() {
	// Dependency would be some database query
}

