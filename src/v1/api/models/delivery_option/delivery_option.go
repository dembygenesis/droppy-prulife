package delivery_option

import "github.com/dembygenesis/droppy-prulife/utilities/database"

func (d *DeliveryOption) ValidateByName() (bool, error) {
	return database.ValidEntry(d.Name, "`name`", "delivery_option")
}