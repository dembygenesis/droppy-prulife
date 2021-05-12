package crud

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/droppy-prulife/src/v1/api/models/deliveries"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/config"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/domain/delivery"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/domain/sysparam"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/domain/user"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/utils"
	"github.com/dembygenesis/droppy-prulife/utilities/database"
	"github.com/dembygenesis/droppy-prulife/utilities/file"
	StringUtility "github.com/dembygenesis/droppy-prulife/utilities/string"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type DeliveryRepository interface {
	Create(p *delivery.RequestCreateDelivery, f *multipart.FileHeader) *utils.ApplicationError
	Update(p *delivery.RequestUpdateDelivery) *utils.ApplicationError
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

// validateSellerId - ensures the seller is valid
func (d *deliveryRepository) validateSellerId(tx *gorm.DB, p *delivery.RequestCreateDelivery) error {
	var _user user.User

	res := tx.Model(&_user).Select("user_type.name AS user_type").
		Joins(`INNER JOIN user_type ON user_type.id = user.user_type_id`).
		Where("user.id = ?", p.SellerId).
		Where("user_type.name = ?", config.UserTypeSeller).
		First(&_user).Error

	if errors.Is(res, gorm.ErrRecordNotFound) {
		return errors.New("seller_id not found")
	}
	if res != nil {
		return errors.New("error validating the seller_id: " + res.Error())
	}
	if _user.UserType != config.UserTypeSeller {
		return errors.New("seller_id must be of type 'Seller'")
	}

	return nil
}

// validateSellerId - ensures the dropshipper is valid
func (d *deliveryRepository) validateDropshipperId(tx *gorm.DB, p *delivery.RequestCreateDelivery) error {
	var _user user.User

	res := tx.Model(&_user).Select("user_type.name AS user_type").
		Joins(`INNER JOIN user_type ON user_type.id = user.user_type_id`).
		Where("user.id = ?", p.DropshipperId).
		Where("user_type.name = ?", config.UserTypeDropshipper).
		First(&_user).Error

	if errors.Is(res, gorm.ErrRecordNotFound) {
		return errors.New("dropshipper_id not found")
	}
	if res != nil {
		return errors.New("error validating the dropshipper_id: " + res.Error())
	}
	if _user.UserType != config.UserTypeDropshipper {
		return errors.New("dropshipper_id must be of type 'Dropshipper'")
	}

	return nil
}

// validateDeliveryOption - ensures the "delivery_option" parameter is valid
func (d *deliveryRepository) validateDeliveryOption(tx *gorm.DB, p *delivery.RequestCreateDelivery) error {
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
func (d *deliveryRepository) validateServiceFeeType(tx *gorm.DB, p *delivery.RequestCreateDelivery) error {
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

func (d *deliveryRepository) getDeliveryStatusName(tx *gorm.DB, deliveryId int) (string, error) {
	// Fetch current delivery status
	var currentDeliveryStatus string

	err := tx.Raw(`
		SELECT 
			ds.name 
		FROM delivery d 
		INNER JOIN delivery_status ds 
			ON 1 = 1
				AND d.delivery_status_id = ds.id 
		WHERE 1 = 1
			AND d.id = ?
	`, deliveryId).Scan(&currentDeliveryStatus).Error

	return currentDeliveryStatus, err
}

// validateDeliveryStatusLogic - validates delivery status transition
func (d *deliveryRepository) validateDeliveryStatusLogic(tx *gorm.DB, p *delivery.RequestUpdateDelivery) error {
	currentDeliveryStatus, err := d.getDeliveryStatusName(tx, p.DeliveryId)
	if err != nil {
		return errors.New("error trying to fetch the current delivery status name: " + err.Error())
	}

	// Logic: Current delivery is already 'Voided'
	if currentDeliveryStatus == config.DeliveryStatusVoided {
		return errors.New("delivery_id is voided, and has no valid actions available")
	}

	// Logic: Delivery item is being 'Voided'
	if p.DeliveryStatus == config.DeliveryStatusVoided {
		// Throw no error and process immediately
		return nil
	}

	// Logic: Delivery item is 'Rejected'
	if currentDeliveryStatus == config.DeliveryStatusRejected {
		return errors.New("delivery_id is rejected, and has no valid actions available - except 'Void'")
	}

	// Logic: Pending Approval
	if currentDeliveryStatus == config.DeliveryStatusVoided {
		return nil
	}

	// Logic: Pending Approval
	if currentDeliveryStatus == config.DeliveryStatusPendingApproval {
		if p.DeliveryStatus != config.DeliveryStatusProposed {
			return errors.New("Pending Approval can only progress to 'Proposed' or 'Voided'")
		}
	}

	// Logic: Proposed
	if currentDeliveryStatus == config.DeliveryStatusProposed {
		if !(p.DeliveryStatus == config.DeliveryStatusAccepted ||
			p.DeliveryStatus == config.DeliveryStatusRejected) {
			return errors.New("Proposed can only progress to 'Accepted' or 'Rejected'")
		}
	}

	// Logic: Accepted
	if currentDeliveryStatus == config.DeliveryStatusAccepted {
		if p.DeliveryStatus != config.DeliveryStatusFulfilled {
			return errors.New("Accepted can only progress to Fulfilled")
		}
	}

	// Logic: Fulfilled
	if currentDeliveryStatus == config.DeliveryStatusFulfilled {
		if !(p.DeliveryStatus == config.DeliveryStatusDelivered ||
			p.DeliveryStatus == config.DeliveryStatusRejected) {
			return errors.New("Fulfilled can only progress to 'Delivered' or 'Rejected'")
		}
	}

	return nil
}

// validateDeliveryStatus - ensures the delivery_status valid
func (d *deliveryRepository) validateDeliveryStatus(tx *gorm.DB, p *delivery.RequestUpdateDelivery) error {
	var count int

	err := tx.Raw("SELECT COUNT(*) FROM delivery_status WHERE name = ?", p.DeliveryStatus).Scan(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("delivery_status does not exist")
	}
	return nil
}

// validateDeliveryId - ensures the seller is valid
func (d *deliveryRepository) validateDeliveryId(tx *gorm.DB, p *delivery.RequestUpdateDelivery) error {
	var count int

	err := tx.Raw("SELECT COUNT(*) FROM delivery WHERE id = ?", p.DeliveryId).Scan(&count).Error
	if err != nil {
		return err
	}
	fmt.Println("count", count)
	if count == 0 {
		return errors.New("delivery_id does not exist")
	}
	return nil
}

// validatePolicyNumber ensures that the "policy_number" to be inserted does not already exist
func (d *deliveryRepository) validatePolicyNumber(tx *gorm.DB, p *delivery.RequestCreateDelivery) error {
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

func (d *deliveryRepository) updateValidations(tx *gorm.DB, p *delivery.RequestUpdateDelivery) error {
	var err error

	// Validate: User Type
	if p.CreatedByUserType != config.UserTypeDropshipper {
		return errors.New("dropshippers are the only one's that can create a new policy")
	}

	// Validate: Delivery ID
	err = d.validateDeliveryId(tx, p)
	if err != nil {
		return err
	}

	// Validate: Delivery Status
	err = d.validateDeliveryStatus(tx, p)
	if err != nil {
		return err
	}

	// Validate: Delivery Transition Status Logic
	err = d.validateDeliveryStatusLogic(tx, p)
	if err != nil {
		return err
	}

	return err
}

func (d *deliveryRepository) createValidations(tx *gorm.DB, p *delivery.RequestCreateDelivery) error {
	var err error

	// Validate: User Type
	if p.CreatedByUserType != config.UserTypeDropshipper {
		return errors.New("dropshippers are the only one's that can create a new policy")
	}

	// Validate: Seller ID
	err = d.validateSellerId(tx, p)
	if err != nil {
		return err
	}

	// Validate: Dropshipper ID
	err = d.validateDropshipperId(tx, p)
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

func (d *deliveryRepository) createDelivery(tx *gorm.DB, p *delivery.RequestCreateDelivery) error {
	var err error
	var serviceFee float64
	var deliveryOptionId int
	var deliveryStatusId int

	// Extract service fee
	var sysParam sysparam.SysParam
	res := tx.Table("sysparam").Select("value").Where("`key` = ?", p.ServiceFeeType).First(&sysParam).Error

	if errors.Is(res, gorm.ErrRecordNotFound) {
		return errors.New("service_fee_type not found")
	}
	if res != nil {
		return errors.New("error fetching the service_fee_type: " + res.Error())
	}
	serviceFee, err = strconv.ParseFloat(sysParam.Value, 64)

	// Extract delivery option id
	var deliveryOption delivery.DeliveryOption
	res = tx.Table("delivery_option").
		Select("`id`").
		Where("name = ?", p.DeliveryOption).
		First(&deliveryOption).Error

	if errors.Is(res, gorm.ErrRecordNotFound) {
		return errors.New("delivery_option not found")
	}
	if res != nil {
		return errors.New("error fetching the delivery option: " + res.Error())
	}
	deliveryOptionId = int(deliveryOption.Id)

	// Extract delivery status id
	err = tx.Table("delivery_status").
		Select("`id`").
		Where("name = ?", config.DeliveryStatusPendingApproval).
		Scan(&deliveryStatusId).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("delivery_option not found")
	}
	if err != nil {
		return errors.New("error fetching the delivery option: " + err.Error())
	}
	if deliveryStatusId == 0 {
		return errors.New("delivery_status_id not found")
	}

	// Insert
	_delivery := map[string]interface{}{
		"created_by":         p.SellerId,
		"updated_by":         p.SellerId,
		"name":               p.Name,
		"address":            p.Address,
		"service_fee":        serviceFee,
		"delivery_option_id": deliveryOptionId,
		"delivery_status_id": deliveryStatusId,
		"seller_id":          p.SellerId,
		"dropshipper_id":     p.DropshipperId,
		"contact_number":     p.ContactNumber,
		"note":               p.Note,
		"image_url":          p.ImageUrl,
		"item_description":   p.ItemDescription,
		"policy_number":      p.PolicyNumber,
	}
	return tx.Model(&delivery.Delivery{}).Create(_delivery).Error
}

func (d *deliveryRepository) uploadToS3AndSync(
	tx *gorm.DB,
	f *multipart.FileHeader,
	p *delivery.RequestCreateDelivery,
) error {
	var err error
	var lastInsertId int

	// Get last insert id
	lastInsertId, err = database.GetLastInsertIDGorm(tx)
	if err != nil {
		return err
	}
	// Send image off
	err = deliveries.UploadDeliveryImageToS3(lastInsertId, f)
	if err != nil {
		return err
	}
	// Sync in database
	bucketPath := os.Getenv("AWS_BUCKET")
	fileType := file.GetMultiPartFileType(f)
	p.ImageUrl = `https://s3-ap-southeast-1.amazonaws.com/` + bucketPath + `/delivery_images/` + strconv.Itoa(lastInsertId) + `_item.` + fileType
	sqlUpdateDeliveryImageUrl := `UPDATE delivery SET image_url = ? WHERE id = ?`
	return tx.Exec(sqlUpdateDeliveryImageUrl, p.ImageUrl, lastInsertId).Error
}

func (d *deliveryRepository) Create(p *delivery.RequestCreateDelivery, f *multipart.FileHeader) *utils.ApplicationError {
	var err error

	err = db.Transaction(func(tx *gorm.DB) error {
		// Validate
		err = d.createValidations(tx, p)
		if err != nil {
			return err
		}

		// Insert
		err = d.createDelivery(tx, p)
		if err != nil {
			return err
		}

		// Upload file to s3 and sync in database
		return d.uploadToS3AndSync(tx, f, p)
	})

	if err != nil {
		return &utils.ApplicationError{
			HttpStatus: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Error:      err,
		}
	}
	return nil
}

// updateDeliveryItemStatus - updates a delivery item's status to another
func (d *deliveryRepository) updateDeliveryItemStatus(tx *gorm.DB, p *delivery.RequestUpdateDelivery) error {
	var err error
	var deliveryStatusId int

	// Extract delivery status
	err = tx.Table("delivery_status").
		Select("`id`").
		Where("name = ?", p.DeliveryStatus).
		Scan(&deliveryStatusId).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("delivery_option not found")
	}
	if err != nil {
		return errors.New("error fetching the delivery option: " + err.Error())
	}
	if deliveryStatusId == 0 {
		return errors.New("delivery_status_id not found")
	}

	sqlUpdateDeliveryImageUrl := `UPDATE delivery SET delivery_status_id = ? WHERE id = ?`
	return tx.Exec(sqlUpdateDeliveryImageUrl, deliveryStatusId, p.DeliveryId).Error
}

func (d *deliveryRepository) updateDelivery(tx *gorm.DB, p *delivery.RequestUpdateDelivery) error {
	var err error

	// TODO: Different transition scenarios here


	// Update delivery item status
	err = d.updateDeliveryItemStatus(tx, p)

	return err
}

func (d *deliveryRepository) Update(p *delivery.RequestUpdateDelivery) *utils.ApplicationError {
	var err error

	err = db.Transaction(func(tx *gorm.DB) error {
		// Validate
		err = d.updateValidations(tx, p)
		if err != nil {
			return err
		}

		// Update
		err = d.updateDelivery(tx, p)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return &utils.ApplicationError{
			HttpStatus: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Error:      err,
		}
	}
	return nil
}