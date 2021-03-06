package services

import (
	"github.com/dembygenesis/droppy-prulife/src/v2/api/domain/delivery"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/repositories/delivery/crud"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/utils"
	"mime/multipart"
	"net/http"
	"reflect"
)

type Service interface {
	// Implement your service methods here
	Create(p *delivery.RequestCreateDelivery, f *multipart.FileHeader) *utils.ApplicationError
	Update(p *delivery.RequestUpdateDelivery) *utils.ApplicationError
}

type service struct {
	deliveryRepo crud.DeliveryRepository
	// Composition of repositories that matches the "Service" interface above

	// Why composition?
}

func NewDeliveryService(deliveryRepo crud.DeliveryRepository) Service {
	return &service{
		deliveryRepo: deliveryRepo,
	}
}

func (s *service) ValidateNoEmptyParamsCreate(p *delivery.RequestCreateDelivery) *utils.ApplicationError {
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
			Message:    "missing parameters",
			Error:      missingParameters,
		}
	} else {
		return nil
	}
}

func (s *service) ValidateNoEmptyParamsUpdate(p *delivery.RequestUpdateDelivery) *utils.ApplicationError {
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
			Message:    "missing parameters",
			Error:      missingParameters,
		}
	} else {
		return nil
	}
}

func (s *service) Create(p *delivery.RequestCreateDelivery, f *multipart.FileHeader) *utils.ApplicationError {
	// Ensure no empty parameters
	appError := s.ValidateNoEmptyParamsCreate(p)
	if appError != nil {
		return appError
	}

	return s.deliveryRepo.Create(p, f)
}

func (s *service) Update(p *delivery.RequestUpdateDelivery) *utils.ApplicationError {
	// Ensure no empty parameters
	appError := s.ValidateNoEmptyParamsUpdate(p)
	if appError != nil {
		return appError
	}

	return s.deliveryRepo.Update(p)
}
