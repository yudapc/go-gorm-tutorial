package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type User struct {
	// gorm.Model
	ID          uint         `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   *time.Time   `json:"deleted_at"`
	Name        string       `json:"name"`
	CreditCards []CreditCard `gorm:"foreignkey:UserRefer" json:"credit_cards"`
}

type CreditCard struct {
	// gorm.Model
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Number    string     `json:"number"`
	UserRefer uint       `json:"user_refer"`
	User      *User      `gorm:"foreignkey:UserRefer" json:"user"`
}

func main() {
	fmt.Println("Running...")

	// Database connection
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	db.LogMode(false) // set to true for see logs gorm
	defer db.Close()

	// Migrate the schema
	migrate(db)

	// Create
	seed(db)

	// Read user with nested credit cards
	getUserById(db, 1)
	getCreditCardByNumber(db, "321")

	fmt.Println("Done!")
}

func migrate(db *gorm.DB) {
	db.DropTableIfExists(&User{}, &CreditCard{})
	db.AutoMigrate(&User{}, &CreditCard{})
}

func seed(db *gorm.DB) {
	user := &User{
		Name: "Jhon", CreditCards: []CreditCard{
			{Number: "123"},
			{Number: "321"},
		},
	}
	db.Create(user)
}

func getUserById(db *gorm.DB, id uint) {
	var user User
	db.Preload("CreditCards").Find(&user, 1)
	response(user, "user detail")
}

func getCreditCardByNumber(db *gorm.DB, number string) {
	var creditCard CreditCard
	db.Preload("User").Find(&creditCard, "number = ?", number)
	response(creditCard, "credit card")
}

func response(data interface{}, info string) {
	fmt.Println("==============================================")
	res, _ := json.Marshal(data)
	fmt.Println(info, ":", string(res))
}
