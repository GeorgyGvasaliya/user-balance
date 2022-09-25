package models

type User struct {
	UserID  int     `json:"user_id"`
	Balance float64 `json:"balance"`
}

type ChangeBalance struct {
	UserID int     `json:"user_id"`
	Money  float64 `json:"money"`
}

type TransferMoney struct {
	FromID int     `json:"from_id"`
	ToID   int     `json:"to_id"`
	Money  float64 `json:"money"`
}

type CurrencyConverter struct {
	Timestamp int                `json:"timestamp"`
	Base      string             `json:"base"`
	Rates     map[string]float64 `json:"rates"`
}
