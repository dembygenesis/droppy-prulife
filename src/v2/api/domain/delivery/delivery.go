package delivery

import "github.com/dembygenesis/droppy-prulife/src/v2/api/repositories/delivery/domain"

// This method must implement validate
type IDelivery interface {
	Validate(i int) bool
}

// This method must have a repository service that implements validate.
// But don't get confused, they don't have to be equl
type Delivery struct {
	Id int
	service domain.DeliveryDomainRepository
}

var (
	delivery IDelivery
)

func init() {
	delivery = &Delivery{}
}

func NewDelivery(i int) IDelivery {
	delivery = &Delivery{Id: i}

	return delivery
}


func (d *Delivery) Validate(i int) bool {
	return d.service.Validate(i)
}