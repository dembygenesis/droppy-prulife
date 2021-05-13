package deliveries

import (
	"errors"
	"fmt"
	"github.com/dembygenesis/droppy-prulife/src/v1/api/database"
	SysParamModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/sysparam"
	UserModel "github.com/dembygenesis/droppy-prulife/src/v1/api/models/users"
	"github.com/dembygenesis/droppy-prulife/utilities/aws/s3"
	"github.com/dembygenesis/droppy-prulife/utilities/file"
	"github.com/jmoiron/sqlx"
	"mime/multipart"
	"os"
	"reflect"
	"regexp"
	"strconv"
)

func InsertCoinTransaction(t *sqlx.Tx, adminAccount int, userId int, creditType string, amount float64, deliveryId int) (error) {
	sql := `
		INSERT INTO coin_transaction (created_by, created_date, is_active, user_id, type, amount, delivery_id)
		VALUES (
			   ?,
			   NOW(),
			   1,
			   ?,
			   ?,
			   ?,
			   ?
		);  
	`

	_, err := t.Exec(sql, adminAccount, userId, creditType, amount, deliveryId)

	return err
}

func UpdateCoinTotals(t *sqlx.Tx, adminId int, userId int, amount float64) error {

	// Check if a totals record already exists
	var totalsCount int
	sqlTotalsExists := `SELECT COUNT(id) FROM user_total WHERE user_id = ?`

	err := t.Get(&totalsCount, sqlTotalsExists, userId)

	if err != nil {
		return errors.New("something went wrong when trying to get the user totals")
	}

	if totalsCount == 0 {
		// Create new totals entry
		sql := `
			INSERT INTO user_total (user_id, amount, coin_amount, created_by, last_updated)
			VALUES (
				   ?,
				   0,
				   ?,
				   ?,
				   NOW()
			);
		`

		_, err := t.Exec(sql, userId, amount, adminId)

		return err
	} else {
		// Update totals
		sql := `
			UPDATE user_total
			SET coin_amount = coin_amount + ?
			WHERE user_id = ?
		`

		_, err := t.Exec(sql, amount, userId)

		return err
	}
}

// IGNORE
func GetImageName(lastInsertDeliveryId int, fileType string) string {
	deliveryId := strconv.Itoa(lastInsertDeliveryId)

	reDb := regexp.MustCompile("image/")
	trimmedFileType := reDb.ReplaceAllString(fileType, "")
	formattedFileName := deliveryId + "_item." + trimmedFileType

	return formattedFileName
}

func UploadDeliveryImageToS3(lastInsertDeliveryId int, f *multipart.FileHeader) error {

	// Conversion
	reader, err := s3.MultipartToReader(f)

	if err != nil {
		return err
	}

	// Content type
	fileType, err := file.GetFileContentType(reader)

	if err != nil {
		return err
	}

	// Build file name
	s3FileName := GetImageName(lastInsertDeliveryId, fileType)
	bucketPath := os.Getenv("AWS_BUCKET")

	// Upload to S3
	// err = s3.UploadObject(s3FileName, reader, fileType, bucketPath + "/delivery_images/")
	err = s3.UploadObjectMultiPart(s3FileName, f, bucketPath + "/delivery_images/")

	if err != nil {
		return err
	}

	return nil
}

func (d *Delivery) ValidateCreate(p ParamsCreateDelivery) []string {
	var missingParameters []string

	s := p
	v := reflect.ValueOf(s)
	v2 := reflect.TypeOf(s)
	typeOfS := v.Type()

	for i := 0; i< v.NumField(); i++ {
		propertyType := typeOfS.Field(i).Type.String()
		propertyValue := v.Field(i).Interface()
		propertyName := v2.Field(i).Tag.Get("json")
		required := v2.Field(i).Tag.Get("required")

		if required != "false" {
			if propertyType == "string" {
				if propertyValue == "" {
					missingParameters = append(missingParameters, `` + propertyName + ` empty`)
				}
			} else if propertyType == "float64" {
				if propertyValue == 0 {
					missingParameters = append(missingParameters, `` + propertyName + ` empty`)
				}
			}
		}
	}

	return missingParameters
}

func (d *Delivery) Create(p *ParamsCreateDelivery, f *multipart.FileHeader, userType string) error {

	// Start transaction
	t := database.DBInstancePublic.MustBegin()

	// Validate user type
	if userType != "Seller" {
		_ = t.Rollback()
		return errors.New("user type must be seller")
	}

	// Validate delivery payment method
	var countDeliveryPaymentMethod int
	sqlValidDeliveryPaymentMethod := `SELECT COUNT(id) FROM delivery_payment_method WHERE name = ?`

	err := t.Get(&countDeliveryPaymentMethod, sqlValidDeliveryPaymentMethod, p.DeliveryPaymentMethod)

	if err != nil {
		_ = t.Rollback()
		return err
	}

	// Validate delivery options
	var countDeliveryOption int
	  sqlValidDeliveryOption := `SELECT COUNT(id) FROM delivery_option WHERE name = ?`

	err = t.Get(&countDeliveryOption, sqlValidDeliveryOption, p.DeliveryOption)

	if err != nil {
		_ = t.Rollback()
		return err
	}

	if countDeliveryOption == 0 {
		_ = t.Rollback()
		return errors.New("no delivery_option found")
	}

	// Guard must be dropship only
	if p.DeliveryOption != "Dropship" {
		_ = t.Rollback()
		return errors.New("delivery option is for dropship only")
	}

	// Validate region
	var countRegion int
	sqlValidRegion := `SELECT COUNT(id) FROM region WHERE name = ?`

	err = t.Get(&countRegion, sqlValidRegion, p.Region)

	if err != nil {
		_ = t.Rollback()
		return err
	}

	if countRegion == 0 {
		_ = t.Rollback()
		return errors.New("no valid region found")
	}

	// Validate declared amount
	if p.DeclaredAmount <= 0 {
		_ = t.Rollback()
		return errors.New("declared_amount must be greater than 0")
	}

	// Validate dropshipper email
	s := SysParamModel.SysParam{}

	sysParam, err := s.GetByKey("HANDLER_DROPSHIP_" + p.Region)

	if err != nil {
		_ = t.Rollback()
		return err
	}

	var validateDropshipper []ValidateDropshipper

	sqlValidateDropshipper := `
		SELECT 
			u.id
		FROM user u 
		WHERE 1 = 1
			AND u.email = ?
			AND u.is_active = 1
			AND u.user_type_id = (SELECT id FROM user_type WHERE name = 'Dropshipper')
	`

	err = t.Select(&validateDropshipper, sqlValidateDropshipper, sysParam.Value)

	if err != nil {
		_ = t.Rollback()
		return err
	}

	if len(validateDropshipper) == 0 {
		_ = t.Rollback()
		return errors.New("No dropshipper avaiable. Check that your user and config match")
	} else {
		p.DropshipperId = validateDropshipper[0].Id
	}

	// Get service fee
	sysParam, err = s.GetByKey("NBA_SERVICE_FEE")

	if err != nil {
		_ = t.Rollback()
		return err
	}

	serviceFee, err := strconv.ParseFloat(sysParam.Value, 64)

	if err != nil {
		_ = t.Rollback()
		return err
	}

	// Check user balances
	var userBalance float64
	sqlUserBalance := `SELECT coin_amount FROM user_total WHERE user_id = ?`

	err = t.Get(&userBalance, sqlUserBalance, p.SellerId)

	if err != nil {
		_ = t.Rollback()
		return errors.New("coin_amount is not available")
	}

	if userBalance < serviceFee {
		_ = t.Rollback()
		return errors.New("user balance is not enough")
	}

	/**
	Do inserts
	 */

	// Insert into delivery
	sqlCreateDelivery := `
		INSERT INTO delivery (
			created_by, created_date, is_active, name, address, region_id,
			service_fee, base_price, declared_amount, delivery_option_id, seller_id,
			dropshipper_id, delivery_status_id, contact_number, note, amount_distributor,
			item_description, delivery_payment_method_id
		)	
		VALUES (
			?,
            NOW(),
            1,
            ?,
            ?,
            (SELECT id FROM region WHERE name = '` + p.Region + `'),
			?,
			?,
            ?,
            (SELECT id FROM delivery_option WHERE name = '` + p.DeliveryOption + `'),
            ?,
            ?,
            (SELECT id FROM delivery_status WHERE name = 'Proposed'),
            ?,
            ?,
            ?,
			?,
			(SELECT id FROM delivery_payment_method WHERE name = '` + p.DeliveryPaymentMethod + `')
		);
	`

	res, err := t.Exec(
		sqlCreateDelivery,
		p.SellerId,
		p.Name,
		p.Address,
		serviceFee,
		0,
		p.DeclaredAmount,
		p.SellerId,
		p.DropshipperId,
		p.ContactNumber,
		p.Note,
		0,
		p.ItemDescription,
	)

	if err != nil {
		_ = t.Rollback()
		return err
	} else {
		fmt.Println("YAY!", res)
	}

	// Get last insert id
	lastInsertDeliveryId, err := database.GetLastInsertID(t)

	if err != nil {
		_ = t.Rollback()
		return err
	}

	// Insert into delivery tracking
	sqlCreateDeliveryTracking := `
		INSERT INTO delivery_tracking (delivery_id, delivery_status_id, last_updated, updated_by)
		VALUES (
				   ?,
				   (SELECT id FROM delivery_status WHERE name = 'Proposed'),
				   NOW(),
				   ?
			   );
	`

	res, err = t.Exec(sqlCreateDeliveryTracking, lastInsertDeliveryId, p.SellerId)

	if err != nil {
		_ = t.Rollback()
		return err
	}

	// Get admin account
	sysParam, err = s.GetByKey("HANDLER_ADMIN")

	if err != nil {
		_ = t.Rollback()
		return err
	}

	adminEmail := sysParam.Value

	u := UserModel.User{Email: adminEmail}
	adminId, err := u.GetByEmail()

	if err != nil {
		_ = t.Rollback()
		return err
	}

	sysParam, err = s.GetByKey("DROPSHIPPER_FEE")

	if err != nil {
		_ = t.Rollback()
		return err
	}

	// Debit
	sellerFee := serviceFee

	// Credit
	dropshipperFee, err := strconv.Atoi(sysParam.Value)

	if err != nil {
		_ = t.Rollback()
		return err
	}


	// MARKER
	// This is the formula for adding a new entry

	// Debit
	adminFee := int(serviceFee) - dropshipperFee

	// Add seller coin transaction
	err = InsertCoinTransaction(t, adminId, p.SellerId, "D", sellerFee, lastInsertDeliveryId)

	if err != nil {
		_ = t.Rollback()
		return err
	}

	// Add dropshipper coin transaction
	err = InsertCoinTransaction(t, adminId, p.DropshipperId, "C", float64(dropshipperFee * -1), lastInsertDeliveryId)

	if err != nil {
		_ = t.Rollback()
		return err
	}

	// Add admin coin transaction
	err = InsertCoinTransaction(t, adminId, adminId, "D", float64(adminFee), lastInsertDeliveryId)

	if err != nil {
		_ = t.Rollback()
		return err
	}

	/**
	Do total updates
	 */

	// Update seller totals
	err = UpdateCoinTotals(t, adminId, p.SellerId, sellerFee * -1)

	if err != nil {
		_ = t.Rollback()
		return err
	}

	// Update dropshipper totals
	err = UpdateCoinTotals(t, adminId, p.DropshipperId, float64(dropshipperFee))

	if err != nil {
		_ = t.Rollback()
		return err
	}

	// Update admin totals
	err = UpdateCoinTotals(t, adminId, adminId, float64(adminFee))

	if err != nil {
		_ = t.Rollback()
		return err
	}

	/**
	Upload to s3
	Either this will transfer an existing S3 file to the new delivery OR
	upload a new file from the client
	*/
	err = UploadDeliveryImageToS3(lastInsertDeliveryId, f)

	if err != nil {
		_ = t.Rollback()
		return err
	}

	// Update image on delivery
	bucketPath := os.Getenv("AWS_BUCKET")
	fileType := file.GetMultiPartFileType(f)
	p.ImageUrl = `https://s3-ap-southeast-1.amazonaws.com/` + bucketPath + `/delivery_images/` + strconv.Itoa(lastInsertDeliveryId) + `_item.` + fileType

	sqlUpdateDeliveryImageUrl := `UPDATE delivery SET image_url = ? WHERE id = ?`
	_, err = t.Exec(sqlUpdateDeliveryImageUrl, p.ImageUrl, lastInsertDeliveryId)

	if err != nil {
		_ = t.Rollback()
		return err
	}

	// Commit transaction
	_ = t.Commit()

	return nil
}
