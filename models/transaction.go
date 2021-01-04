package models

import "time"

type TransactionType string

const (
	income = iota
	expense
)

type Transaction struct {
	CommonFields
	Type        TransactionType `json:"type"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date"`
	Amount      int             `json:"amount"`
	CategoryID  int             `json:"category_id"`
	UserID      int             `json:"user_id"`
}
