package utils

// EncloseString - encloses strings on in a delimiter
func EncloseString(s string, d string) string {
	return ``+ d +``+ s +``+ d +``
}

// GetSQLValue - returns a column with the same value as itself if the string is empty
func GetSQLValue(column string, value string) string {
	if value == "" {
		return EncloseString(column, "`")
	}
	return EncloseString(value, "\"")
}