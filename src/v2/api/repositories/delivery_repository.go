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


func (d *deliveryRepository) validateSellerId(tx *gorm.DB, p *delivery.CreateDelivery) error {
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

func (d *deliveryRepository) validateDeliveryOption(tx *gorm.DB, p *delivery.CreateDelivery) error {
	var deliveryOption delivery.DeliveryOption

	res := tx.Table("delivery_option").
		Select("name").
		Where("name = ?", p.DeliveryOption).
		First(&deliveryOption).Error

	// No record found
	if errors.Is(res, gorm.ErrRecordNotFound) {
		return errors.New("delivery_option not found")
	}
	// MYSQL err
	if res != nil {
		return errors.New("error validating the delivery option: " + res.Error())
	}
	// Dropship only...
	if deliveryOption.Name != "Dropship" {
		return errors.New("delivery_option must be dropship only: ")
	}
	return nil
}

func (d *deliveryRepository) validations(tx *gorm.DB, p *delivery.CreateDelivery) error {
	var err error

	// Validate: Seller ID
	err = d.validateSellerId(tx, p)
	if err != nil {
		return err
	}

	// Validate: Delivery Option
	err = d.validateDeliveryOption(tx, p)
	if err != nil {
		return err
	}

	// Validate:

	return nil
}

func (d *deliveryRepository) Create(p *delivery.CreateDelivery) *utils.ApplicationError {

	var err error

	err = db.Transaction(func(tx *gorm.DB) error {
		// Validation logic
		err = d.validations(tx, p)
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
