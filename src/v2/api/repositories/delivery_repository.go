package repositories

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/domain/delivery"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/domain/user"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/utils"
	"github.com/dembygenesis/droppy-prulife/utilities/database"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type DeliveryRepository interface {
	Create(p *delivery.CreateDelivery) *utils.ApplicationError
}

type deliveryRepository struct {
}

var db *gorm.DB

func init() {
	// Parse ENV
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get Gorm instance
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDatabase := os.Getenv("DB_DATABASE")
	dbPort := os.Getenv("DB_PORT")

	db, err = database.GetGormInstance(
		dbHost,
		dbUser,
		dbPassword,
		dbDatabase,
		dbPort,
	)

	if err != nil {
		panic("failed to connect database")
	}
}

// Constructor
func NewDeliveryRepository() DeliveryRepository {
	return &deliveryRepository{}
}


func validateSellerId(tx *gorm.DB, p *delivery.CreateDelivery) error {
	var _user user.User

	res := tx.Model(&_user).Select("user_type.name AS user_type").
		Joins(`INNER JOIN user_type ON user_type.id = user.user_type_id`).
		Where("user.id = ?", p.SellerId).
		First(&_user).Error

	// No record found
	if errors.Is(res, gorm.ErrRecordNotFound) {
		return errors.New("seller_id not found")
	}
	// MYSQL err
	if res != nil {
		return errors.New("error validating the user id: " + res.Error())
	}
	// Non-seller error
	if _user.UserType != "Seller" {
		fmt.Println("I fail")
		return errors.New("seller_id must be of type 'Seller'")
	}

	return nil
}

func validations(tx *gorm.DB, p *delivery.CreateDelivery) error {
	var err error

	// Validate: Seller ID
	err = validateSellerId(tx, p)
	if err != nil {
		return err
	}

	return nil
}

func (d *deliveryRepository) Create(p *delivery.CreateDelivery) *utils.ApplicationError {

	var err error

	err = db.Transaction(func(tx *gorm.DB) error {
		// Validation logic
		err = validations(tx, p)
		if err != nil {
			return err
		}

		// Process insert

		return nil
	})

	if err != nil {
		return &utils.ApplicationError{
			HttpStatus: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Error:      err,
			Data:       nil,
		}
	}


	// Validate user types

	// Validate coin entries

	// Perform insert

	// Perform coin updates


	fmt.Println("dan pena")
	// Ok... Let's do the validation's here VIA SQLX


	return nil
}
