package models

import (
	"time"
)

// User struct for holding user data
type User struct {
	ID           int				`json:"id"`
	FirstName    string				`json:"first_name"`
	LastName     string				`json:"last_name"`
	DateOfBirth  time.Time			`json:"date_of_birth"`
	Avatar 		 string				`json:"avatar"`
	BankCardInfo map[int]*BankCard	`json:"_"`
}

// BankCard struct for holding User Bank Account Data
type BankCard struct {
	ID         int		`json:"id"`
	UserID 	   int		`json:"user_id"`
	BankName   string	`json:"bank_name"`
	CardNumber string	`json:"card_number"`
}