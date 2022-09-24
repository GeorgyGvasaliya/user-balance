package models

type User struct {
	UserID  int     `json:"user_id"`
	Balance float64 `json:"balance"`
}

type ChangeBalance struct {
	UserID int     `json:"user_id"`
	Money  float64 `json:"money"`
}
