package services

import (
	"github.com/dembygenesis/droppy-prulife/src/v2/api/domain/delivery"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/repositories"
	"github.com/dembygenesis/droppy-prulife/src/v2/api/utils"
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

func (s *service) Create(p *delivery.CreateDelivery) *utils.ApplicationError {
	return s.deliveryRepo.Create(p)
}
