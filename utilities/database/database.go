package database

import (
	"fmt"
	"github.com/dembygenesis/droppy-prulife/src/v1/api/database"
)

func ValidEntry(v string, c string, t string) (bool, error) {
	var count int
	var hasEntry bool

	sql := `
		SELECT COUNT(*) FROM ` + t + `
		WHERE 1 = 1
			AND ` + c + ` = ?
	`

	err := database.DBInstancePublic.Get(&count, sql, v)

	if err != nil {
		return hasEntry, err
	}

	if count > 0 {
		hasEntry = true
	} else {
		hasEntry = false
	}

	fmt.Println(sql, "count", count, "hasEntry", hasEntry)

	return hasEntry, err
}
