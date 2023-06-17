package main

import (
	model "orm-golang/model"
	service "orm-golang/service"
)

var DB = service.GetDBConnection()

func main() {
	DB.AutoMigrate(&model.Language{})
	DB.AutoMigrate(&model.Category{})
	DB.AutoMigrate(&model.Film{})
	DB.Migrator().CreateConstraint(&model.Film{}, "Language")

	DB.AutoMigrate(&model.FilmCategory{})
	DB.Migrator().CreateConstraint(&model.FilmCategory{}, "Film")
	DB.Migrator().CreateConstraint(&model.FilmCategory{}, "Category")

	DB.AutoMigrate(&model.Actor{})
	DB.AutoMigrate(&model.FilmActor{})
	DB.Migrator().CreateConstraint(&model.FilmActor{}, "Actor")
	DB.Migrator().CreateConstraint(&model.FilmActor{}, "Film")

	DB.AutoMigrate(&model.Country{})
	DB.AutoMigrate(&model.City{})
	DB.Migrator().CreateConstraint(&model.City{}, "Country")

	DB.AutoMigrate(&model.Address{})
	DB.Migrator().CreateConstraint(&model.City{}, "City")

	DB.AutoMigrate(&model.Store{})
	DB.Migrator().CreateConstraint(&model.Store{}, "Address")

	DB.AutoMigrate(&model.Staff{})
	DB.Migrator().CreateConstraint(&model.Staff{}, "Address")
	DB.Migrator().CreateConstraint(&model.Staff{}, "Store")

	//SETING UP THE FOREIGN KEY AFTER CREATING STAFF
	DB.Migrator().CreateConstraint(&model.Store{}, "Staff")

	DB.AutoMigrate(&model.Customer{})
	DB.Migrator().CreateConstraint(&model.Customer{}, "Store")
	DB.Migrator().CreateConstraint(&model.Customer{}, "Address")

	DB.AutoMigrate(&model.Inventory{})
	DB.Migrator().CreateConstraint(&model.Inventory{}, "Film")
	DB.Migrator().CreateConstraint(&model.Inventory{}, "Store")

	DB.AutoMigrate(&model.Rental{})
	DB.Migrator().CreateConstraint(&model.Rental{}, "Inventory")
	DB.Migrator().CreateConstraint(&model.Rental{}, "Customer")
	DB.Migrator().CreateConstraint(&model.Rental{}, "Staff")

	DB.AutoMigrate(&model.Payment{})
	DB.Migrator().CreateConstraint(&model.Payment{}, "Customer")
	DB.Migrator().CreateConstraint(&model.Payment{}, "Staff")
	DB.Migrator().CreateConstraint(&model.Payment{}, "Rental")
}
