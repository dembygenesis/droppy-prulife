package repositories

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/config"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/domain/delivery"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/domain/sysparam"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/domain/user"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/utils"
	"github.com/dembygenesis/droppy-prulife/utilities/database"
	StringUtility "github.com/dembygenesis/droppy-prulife/utilities/string"
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

// validateSellerId - ensures the seller
func (d *deliveryRepository) validateSellerId(tx *gorm.DB, p *delivery.CreateDelivery) error {
	var _user user.User

	res := tx.Model(&_user).Select("user_type.name AS user_type").
		Joins(`INNER JOIN user_type ON user_type.id = user.user_type_id`).
		Where("user.id = ?", p.SellerId).
		Where("user_type.name = ?", config.UserTypeSeller).
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
		return errors.New("seller_id must be of type 'Seller'")
	}

	return nil
}

// validateDeliveryOption - ensures the "delivery_option" parameter is valid
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

// validateServiceFeeType - ensures the "service_fee_type" exists in the "sysparam" and table
// and has values in the "config.ValidServiceFeeTypes" (array of strings)
func (d *deliveryRepository) validateServiceFeeType(tx *gorm.DB, p *delivery.CreateDelivery) error {
	// Validate Service Fee Type
	if StringUtility.StringInSlice(p.ServiceFeeType, config.ValidServiceFeeTypes) == false {
		return errors.New("service_fee_type is invalid")
	}

	// Fetch service fee amount
	var sysParam sysparam.SysParam
	res := tx.Table("sysparam").Select("value").Where("`key` = ?", p.ServiceFeeType).First(&sysParam).Error

	// No record found
	if errors.Is(res, gorm.ErrRecordNotFound) {
		return errors.New("service_fee_type not found")
	}
	// MYSQL err
	if res != nil {
		return errors.New("error validating the service_fee_type: " + res.Error())
	}

	// Fetch User balance
	/*var _user user.User

	res = tx.Table("user_total").
		Select("coin_amount AS user_balance").
		Where("user_id = ?", p.SellerId).
		First(&_user).Error

	// No record found
	if errors.Is(res, gorm.ErrRecordNotFound) {
		return errors.New("delivery_option not found")
	}
	// MYSQL err
	if res != nil {
		return errors.New("error validating the delivery option: " + res.Error())
	}

	// We don't need this because
	userBalance, _ := strconv.ParseFloat(sysParam.Value,64)
	if userBalance > _user.UserBalance {
		return errors.New("user cannot afford the transaction")
	}*/

	return nil
}

// validatePolicyNumber ensures that the "policy_number" to be inserted does not already exist
func (d *deliveryRepository) validatePolicyNumber(tx *gorm.DB, p *delivery.CreateDelivery) error {
	var count int

	err := tx.Raw("SELECT COUNT(*) FROM delivery WHERE policy_number = ?", p.PolicyNumber).Scan(&count).Error
	if err != nil {
		return err
	}
	if count == 1 {
		return errors.New("policy_number already exists")
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

	// Validate: Service Fee Type
	err = d.validateServiceFeeType(tx, p)
	if err != nil {
		return err
	}

	// Validate: Policy Number
	err = d.validatePolicyNumber(tx, p)
	if err != nil {
		return err
	}

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

		// No, we don't need this step for now, because this will on the "Update" command
		// Process coin updates

		return nil
	})

	if err != nil {
		return &utils.ApplicationError{
			HttpStatus: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Error:      err,
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
