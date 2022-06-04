package model

type Transaction struct {
	Id        int    `json:"transaction_id" db:"id"`
	UserId    int    `json:"user_id" db:"user_id"`
	Amount    int64  `json:"amount" db:"amount"`
	Operation string `json:"operation" db:"operation"`
	Date      string `json:"date" db:"date"`
}

type TransctionInfo struct {
}
