package controller

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var connection *gorm.DB = connectDB()

func connectDB() *gorm.DB {
	var dbUser string = os.Getenv("TEST_DB_USER")
	var dbPass string = os.Getenv("TEST_DB_PASSWORD")
	var dbName string = os.Getenv("TEST_DB_NAME")
	var dbHost string = os.Getenv("TEST_DB_HOST")
	var dbPort string = os.Getenv("TEST_DB_PORT")

	fmt.Println("Initilizing database connection for tests")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable Timezone=America/Sao_Paulo", dbHost, dbUser, dbPass, dbName, dbPort)
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		fmt.Errorf("Could not connect to test database.")
		panic(err)
	}
	fmt.Println("Connection was successful")
	return connection
}

func GetDBTestConnection() *gorm.DB {
	sqlDB, err := connection.DB()
	if err != nil {
		panic(err)
	}
	newConn, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return newConn
}
