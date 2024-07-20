package models

import (
	"github.com/uadmin/uadmin"
)

type MenuItems struct {
	uadmin.Model
	MenuItems string `uadmin:"file"`
}

func (m *MenuItems) Save() {
	uadmin.Save(m)
	uadmin.Trail(uadmin.INFO, "Saved. %s", m.MenuItems)
	// path := m.MenuItems[1:]
	// f, err := excelize.OpenFile(path)
	// if err != nil {
	// 	uadmin.Trail(uadmin.ERROR, "Unable to open excel file. %s", err)
	// 	return
	// }

	// menuList := f.GetSheetList()

	// for _, menuSheet := range menuList {
	// 	uadmin.Trail(uadmin.INFO, "MENU: %s", menuSheet)
	// 	rows, err := f.GetRows(menuSheet)
	// 	if err != nil {
	// 		uadmin.Trail(uadmin.ERROR, "Unable to get rows. %s", err)
	// 		return
	// 	}

	// 	for _, row := range rows {
	// 		menu := Menus{}
	// 		counter := 0
	// 		m := ""

	// 		for _, colCell := range row {
	// 			colCell = strings.ReplaceAll(colCell, "\n", " ")
	// 			if counter == 0 {
	// 				m += colCell + ";"
	// 				counter++
	// 			} else {
	// 				if counter%11 != 0 {
	// 					m += colCell + ";"
	// 					counter++
	// 				}
	// 				if counter%11 == 0 {
	// 					m += colCell + ";"
	// 					counter = 0

	// 					if uadmin.Count(&menu, "menu = ? AND description = ? AND category = ? AND price = ? AND abbreviation = ? AND from = ? AND to = ? AND image = ? AND discount = ? AND system_code = ? AND kitchen_location = ?", strings.Split(m, ";")[0], strings.Split(m, ";")[1], strings.Split(m, ";")[2], strings.Split(m, ";")[3], strings.Split(m, ";")[4], strings.Split(m, ";")[5], strings.Split(m, ";")[6], strings.Split(m, ";")[7], strings.Split(m, ";")[8], strings.Split(m, ";")[9], strings.Split(m, ";")[10]) == 0 {
	// 						menu.Name = strings.Split(m, ";")[0]
	// 						menu.Description = strings.Split(m, ";")[1]
	// 						menu.Category = strings.Split(m, ";")[2]
	// 						menu.Price = strings.Split(m, ";")[3]
	// 						menu.Abbreviation = strings.Split(m, ";")[4]
	// 						menu.From = strings.Split(m, ";")[5]
	// 						menu.To = strings.Split(m, ";")[6]
	// 						menu.Image = strings.Split(m, ";")[7]
	// 						menu.Discount = strings.Split(m, ";")[8]
	// 						menu.SystemCode = strings.Split(m, ";")[9]
	// 						menu.KitchenLocation = strings.Split(m, ";")[10]

	// 						uadmin.Trail(uadmin.DEBUG, "Menu Name saved: %s", strings.Split(m, ";")[0])
	// 						uadmin.Save(&menu)
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }
}
