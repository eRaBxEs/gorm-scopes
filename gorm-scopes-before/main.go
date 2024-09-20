package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	gorm.Model
	Name   string
	Email  string
	Orders []Order
}

type Order struct {
	UserId      int
	OrderTime   time.Time
	PaymentMode string // cash or card
	Price       int
	User        User
}

var DB *gorm.DB

func connectDatabase() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	dsn := "henry_dev:devdba_user@tcp(127.0.0.1:3307)/gorm_testdb?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})

	if err != nil {
		panic("Failed to connect to databse!")
	}

	DB = database
}

func dbMigrate() {
	DB.AutoMigrate(&User{}, &Order{})
}

func CardOrders(db *gorm.DB) *gorm.DB { // returning this resusable piece of query in scopes
	return db.Where("payment_mode = ?", "card")
}

func PriceGreaterThan30(db *gorm.DB) *gorm.DB {
	return db.Where("price > ?", 30)
}

func UsersFromDomain(domain string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("email like ?", "%"+domain)
	}
}

func main() {

	connectDatabase()
	dbMigrate()

	var orders []Order
	DB.Scopes(CardOrders, PriceGreaterThan30).Find(&orders) // add another argument to the scope function

	// print the details
	for _, order := range orders {
		fmt.Printf("Price: %d, Payment type: %s\n", order.Price, order.PaymentMode)
	}

	var users []User
	DB.Scopes(UsersFromDomain("gmail.com")).Find(&users)

	fmt.Printf("Users:\n")
	for _, user := range users {
		fmt.Printf("User email: %s\n", user.Email)
	}

}
