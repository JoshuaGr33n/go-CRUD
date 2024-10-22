package main

import (
	// "database/sql"
	// _ "github.com/go-sql-driver/mysql"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// var db *sql.DB
var db *gorm.DB


type User struct {
    gorm.Model
    Username string `gorm:"unique;not null"`
    Email    string `gorm:"unique;not null"`
    Password string `gorm:"not null"`
}
// func getMySQLDB() *sql.DB {
// 	// db, err := sql.Open("mysql", "root:@tcp(localhost:8889)/go_crud")
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// // defer db.Close()
// 	// return db
	

// }

func getMySQLDB() *gorm.DB {
	dsn := "root:@tcp(localhost:8889)/go_crud?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db

}

func (Crud) TableName() string {
	return "crud"
}

// Automatically create or migrate tables
func migrateDB() {
	if err := db.AutoMigrate(&Crud{}); err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}

	err := db.AutoMigrate(&User{})
    if err != nil {
        panic("failed to migrate database")
    }
}