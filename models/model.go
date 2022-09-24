package models

type User struct {
	UserID  int `json:"user_id"`
	Balance int `json:"balance"`
}

type AddMoney struct {
	UserID int `json:"user_id"`
	Money  int `json:"money"`
}
