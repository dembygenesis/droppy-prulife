package services

import (
	"errors"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/domain/delivery"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/repositories"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/utils"
	"net/http"
	"reflect"
)

type Service interface {
	// Implement your service methods here
	Create(p *delivery.CreateDelivery) *utils.ApplicationError
}

type service struct {
	deliveryRepo repositories.DeliveryRepository
	// Composition of repositories that matches the "Service" interface above

	// Why composition?

}

func NewDeliveryService(deliveryRepo repositories.DeliveryRepository) Service {
	return &service{
		deliveryRepo: deliveryRepo,
	}
}

func (s *service) ValidateNoEmptyParams(p *delivery.CreateDelivery) *utils.ApplicationError {
	var missingParameters []string

	pp := *p
	v := reflect.ValueOf(pp)
	v2 := reflect.TypeOf(pp)

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
	if len(missingParameters) > 0 {
		return &utils.ApplicationError{
			HttpStatus: http.StatusUnprocessableEntity,
			Message:    "validation errors",
			Error:      errors.New("missing parameters"),
			Data:       missingParameters,
		}
	} else {
		return nil
	}
}

func (s *service) Create(p *delivery.CreateDelivery) *utils.ApplicationError {
	// Ensure no empty parameters
	appError := s.ValidateNoEmptyParams(p)

	if appError != nil {
		return appError
	}

	return s.deliveryRepo.Create(p)
}
