package database

import (
	"fmt"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDb() *gorm.DB { // OOP Constructor
	Db = connectDB()
	return Db
}

func connectDB() *gorm.DB {
	// Connect to SQL Server sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm
	dsn := "sqlserver://admin1234:admin1234@localhost:1433?database=GolangDB"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Error...")
		return nil
	}
	return db
}
