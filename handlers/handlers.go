package handlers

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"user-balance/database"
	"user-balance/models"
)

type Handler struct {
	cfg database.Config
}

func NewHandler(cfg database.Config) Handler {
	return Handler{
		cfg: cfg,
	}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.GET("/get", h.GetUser)
	router.POST("/add", h.AddMoney)
	router.POST("/withdraw", h.Withdraw)
	router.POST("/send", h.SendMoney)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := r.URL.Query().Get("id")
	db := database.NewPostgresDB(h.cfg)
	defer db.Close()

	query := "select * from public.users where user_id=$1"
	var res models.User
	err := db.QueryRow(query, userID).Scan(&res.UserID, &res.Balance)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("There is no user with such ID or wrong query"))
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AddMoney(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	db := database.NewPostgresDB(h.cfg)
	defer db.Close()
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var ch models.ChangeBalance
	err := json.Unmarshal(body, &ch)
	if err != nil {
		log.Println("Cannot Unmarshal data")
		return
	}

	query := "INSERT INTO public.users (user_id, balance) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET balance=users.balance+$2;"
	db.QueryRow(query, ch.UserID, ch.Money)
}

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	db := database.NewPostgresDB(h.cfg)
	defer db.Close()
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var ch models.ChangeBalance
	err := json.Unmarshal(body, &ch)
	if err != nil {
		log.Println("Cannot Unmarshal data")
		return
	}

	queryGet := "select * from public.users where user_id=$1"
	var res models.User
	err = db.QueryRow(queryGet, ch.UserID).Scan(&res.UserID, &res.Balance)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("There is no user with such ID or wrong query"))
		return
	}

	newBalance := res.Balance - ch.Money
	if newBalance < 0 {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("You have no money to purchase"))
		return
	}

	query := "UPDATE public.users SET balance=balance-$1 WHERE user_id=$2;"
	db.QueryRow(query, ch.Money, ch.UserID)
}

func (h *Handler) SendMoney(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	panic("todo")
}
