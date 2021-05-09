package database

// DatabaseCredentials struct for parsing our json.
type DatabaseCredentials struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"Database"`
}