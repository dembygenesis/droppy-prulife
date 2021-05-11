package database

import (
	"fmt"
	"github.com/dembygenesis/droppy-prulife/src/v1/api/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
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

// GetGormInstance - returns a gorm instance
func GetGormInstance(
	dbHost string,
	dbUser string,
	dbPassword string,
	dbDatabase string,
	dbPort string,
) (*gorm.DB, error) {
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbDatabase + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	return db, err
}

func GetLastInsertIDGorm(tx *gorm.DB) (int, error) {
	var lastInsertId int
	err := tx.Raw(`SELECT LAST_INSERT_ID()`).Scan(&lastInsertId).Error
	return lastInsertId, err
}