package model

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var err error
	var dbUser string = os.Getenv("DB_USER")
	var dbPass string = os.Getenv("DB_PASSWORD")
	var dbName string = os.Getenv("DB_NAME")
	var dbHost string = os.Getenv("DB_HOST")

	fmt.Println("Inicializando a conexão com banco de dados")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5433 sslmode=disable Timezone=America/Sao_Paulo", dbHost, dbUser, dbPass, dbName)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Conexão realizada com sucesso")
}
