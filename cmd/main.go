package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"user-balance/handlers"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/get", handlers.GetUser)
	r.HandleFunc("/add", handlers.AddMoney)
	r.HandleFunc("/withdraw", handlers.Withdraw)
	r.HandleFunc("/send", handlers.SendMoney)

	http.ListenAndServe(":8000", r)
}
