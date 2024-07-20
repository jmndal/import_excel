package models

import "github.com/uadmin/uadmin"

type Menus struct {
	uadmin.Model
	Name            string
	Description     string
	Category        string
	Price           float64
	Abbreviation    string
	From            string
	To              string
	Image           string `uadmin:"image"`
	Discount        string
	SystemCode      string
	KitchenLocation string
}

func (m *Menus) String() string {
	return m.Name
}
