package db

import (
	"fmt"
	"log"
	"myproject/pkg/config"

	"database/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *sql.DB

func ConnectPGDB(cnf config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Jakarta", cnf.PGHost, cnf.PGUserName, cnf.PGPassword, cnf.PGDBName, cnf.PgPort)
	fmt.Println("this is the database ", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get *sql.DB: %w", err)
	}

	log.Println("connected to postgres db successfully!")
	DB = sqlDB
	return sqlDB, nil
}
