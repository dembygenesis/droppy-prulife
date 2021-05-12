package domain

// Declare interface
type DeliveryDomainRepository interface {
	Validate(i int) bool
}