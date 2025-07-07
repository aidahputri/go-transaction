package model

type Transaction struct {
	FromAccount string `json:"fromAccount"`
	ToAccount string `json:"toAccount"`
	Amount float64 `json:"amount"`
}