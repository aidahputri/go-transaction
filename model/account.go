package model

type Account struct {
	Name string
	AccountNumber string
	Balance float64
	Blacklisted bool
	Underwatch bool
}