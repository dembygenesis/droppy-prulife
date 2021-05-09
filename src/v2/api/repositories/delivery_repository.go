package repositories

import (
	"errors"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/domain/delivery"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/utils"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"reflect"
)

type DeliveryRepository interface {
	Create(p *delivery.CreateDelivery) *utils.ApplicationError
}

type deliveryRepository struct {
}

var db *gorm.DB

func init() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDatabase := os.Getenv("DB_DATABASE")
	dbPort := os.Getenv("DB_PORT")

	// Initialize gorm
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbDatabase + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}

func NewDeliveryRepository() DeliveryRepository {
	return &deliveryRepository{}
}

// ========================
// BankType
type BankType struct {
	Id   uint `json:"id"`
	Name string
}

func (BankType) TableName() string {
	return "bank_type"
}

// ========================

// Validate validates all parameters passed before proceeding to the insert
func (d *deliveryRepository) Validate(p *delivery.CreateDelivery) []string {
	var missingParameters []string

	s := *p
	v := reflect.ValueOf(s)
	v2 := reflect.TypeOf(s)

	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		propertyType := typeOfS.Field(i).Type.String()
		propertyValue := v.Field(i).Interface()
		propertyName := v2.Field(i).Tag.Get("json")
		required := v2.Field(i).Tag.Get("required")

		if required != "false" {
			if propertyType == "string" {
				if propertyValue == "" {
					missingParameters = append(missingParameters, ``+propertyName+` empty`)
				}
			} else if propertyType == "float64" {
				if propertyValue == 0 {
					missingParameters = append(missingParameters, ``+propertyName+` empty`)
				}
			}
		}
	}

	return missingParameters
}

func (d *deliveryRepository) Create(p *delivery.CreateDelivery) *utils.ApplicationError {
	// Validation
	missingParameters := d.Validate(p)

	if len(missingParameters) > 0 {
		return &utils.ApplicationError{
			HttpStatus: http.StatusUnprocessableEntity,
			Message:    "validation errors",
			Error:      errors.New("missing_parameters"),
			Data:       missingParameters,
		}
	}

	return nil
}
