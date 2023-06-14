package main

import (
	model"orm-golang/model"
)

func main() {
	model.DB.AutoMigrate(&model.Language{})
	model.DB.AutoMigrate(&model.Category{})
	model.DB.AutoMigrate(&model.Film{})
	model.DB.Migrator().CreateConstraint(&model.Film{}, "Language")

	model.DB.AutoMigrate(&model.FilmCategory{})
	model.DB.Migrator().CreateConstraint(&model.FilmCategory{}, "Film")
	model.DB.Migrator().CreateConstraint(&model.FilmCategory{}, "Category")

	model.DB.AutoMigrate(&model.Actor{})
	model.DB.AutoMigrate(&model.FilmActor{})
	model.DB.Migrator().CreateConstraint(&model.FilmActor{}, "Actor")
	model.DB.Migrator().CreateConstraint(&model.FilmActor{}, "Film")

	model.DB.AutoMigrate(&model.Country{})
	model.DB.AutoMigrate(&model.City{})
	model.DB.Migrator().CreateConstraint(&model.City{}, "Country")

	model.DB.AutoMigrate(&model.Address{})
	model.DB.Migrator().CreateConstraint(&model.City{}, "City")

	model.DB.AutoMigrate(&model.Store{})
	model.DB.Migrator().CreateConstraint(&model.Store{}, "Address")

	model.DB.AutoMigrate(&model.Staff{})
	model.DB.Migrator().CreateConstraint(&model.Staff{}, "Address")
	model.DB.Migrator().CreateConstraint(&model.Staff{}, "Store")
	
	//SETING UP THE FOREIGN KEY AFTER CREATING STAFF
	model.DB.Migrator().CreateConstraint(&model.Store{}, "Staff")
	
	model.DB.AutoMigrate(&model.Customer{})
	model.DB.Migrator().CreateConstraint(&model.Customer{},"Store")
	model.DB.Migrator().CreateConstraint(&model.Customer{},"Address")

	model.DB.AutoMigrate(&model.Inventory{})
	model.DB.Migrator().CreateConstraint(&model.Inventory{}, "Film")
	model.DB.Migrator().CreateConstraint(&model.Inventory{}, "Store")

	model.DB.AutoMigrate(&model.Rental{})
	model.DB.Migrator().CreateConstraint(&model.Rental{}, "Inventory")
	model.DB.Migrator().CreateConstraint(&model.Rental{}, "Customer")
	model.DB.Migrator().CreateConstraint(&model.Rental{}, "Staff")

	model.DB.AutoMigrate(&model.Payment{})
	model.DB.Migrator().CreateConstraint(&model.Payment{}, "Customer")
	model.DB.Migrator().CreateConstraint(&model.Payment{}, "Staff")
	model.DB.Migrator().CreateConstraint(&model.Payment{}, "Rental")
}