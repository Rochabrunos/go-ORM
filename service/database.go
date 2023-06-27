package service

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBService struct {
	db *gorm.DB
}

var conn DBService

func init() {
	var err error
	var dbHost string = os.Getenv("DB_HOST")
	var dbUser string = os.Getenv("DB_USER")
	var dbPass string = os.Getenv("DB_PASSWORD")
	var dbName string = os.Getenv("DB_NAME")
	var dbPort string = os.Getenv("DB_PORT")

	fmt.Println("Inicializando a conexão com banco de dados")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable Timezone=America/Sao_Paulo",
		dbHost, dbUser, dbPass, dbName, dbPort)
	conn.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Conexão realizada com sucesso")
}

func GetDBConnection() *gorm.DB {
	sqlDB, err := conn.db.DB()
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
