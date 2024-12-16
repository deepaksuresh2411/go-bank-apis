package main

import (
	"math/rand"
	"time"
)

type CreateAccountReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Account struct {
	ID            int       `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	AccountNumber int       `json:"account_number"`
	Balance       int       `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
}

func NewAccount(fname, lname string) *Account {
	return &Account{
		ID:            rand.Intn(10000),
		FirstName:     fname,
		LastName:      lname,
		AccountNumber: rand.Intn(100000),
		CreatedAt:     time.Now().UTC(),
	}
}
