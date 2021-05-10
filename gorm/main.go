package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type BankType struct {
	Id   uint `gorm:"primaryKey"`
	Name string
}

var db *gorm.DB

func init() {
	var err error
	dsn := "root:root@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: newLogger,
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("failed to connect database")
	}
}

func useTxOne(tx *gorm.DB) error {

	// Joins

	// Declare table
	var users []User

	result := tx.Model(&User{}).Select("user.firstname").Joins(`
		INNER JOIN user_type ON user.user_type_id = user_type.id 
	`).Scan(&users)

	if result.Error != nil {

	} else {
		fmt.Println("users", users)
	}

	return nil
}

func useTxTwo() {

}

func useTransactions() error {
	// Insert two entries to a phonebook VIA transactions

	// How to start transaction, and force an error...

	// Tip: use RAW queries

	// Also, pass the tx object to a function

	var err error

	// Start tx by initializing GORM with a certain config...
	err = db.Transaction(func(tx *gorm.DB) error {
		// Do some query...
		/*var bankType []BankType
		result := tx.Find(&bankType)
		if result.Error != nil {
			return result.Error
		} else {
			fmt.Println("bankType", bankType)
		}*/

		// Do some inserts...
		/*err = tx.Create(&BankType{Name: "Angely"}).Error

		if err != nil {
			return err
		}*/

		fmt.Println()
		fmt.Println()

		useTxOne(tx)

		return nil
	})

	if err != nil {
		fmt.Println("you are DEAD", err.Error())
	} else {
		fmt.Println("successfully committed")
	}

	return err
}

func main() {
	err := useTransactions()

	if err != nil {
		fmt.Println("useTransactions error")
	} else {
		fmt.Println("useTransactions")
	}
}
