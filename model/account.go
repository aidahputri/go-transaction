package model

type Account struct {
	Name string `json:"name"`
	AccountNumber string `json:"accountNumber"`
	Balance float64 `json:"balance"`
	Blacklisted bool `json:"blacklisted"`
	Underwatch bool `json:"underwatch"`
}