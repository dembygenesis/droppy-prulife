package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type BankType struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
}

func (BankType) TableName() string {
	return "bank_type"
}

var db *gorm.DB

func init() {
	var err error
	dsn := "root:root@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}

func main() {
	var bankType []BankType

	result := db.Find(&bankType)

	if result.Error != nil {
		fmt.Println("There is an error when performing the query: " + result.Error.Error())
	} else {
		fmt.Println(bankType)
	}
}
