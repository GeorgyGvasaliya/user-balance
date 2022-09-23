package handlers

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"user-balance/database"
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
	data := database.NewPostgresDB(h.cfg)
	userID := r.URL.Query().Get("id")

	var balance int
	queryGet := "select balance from public.wallet where user_id=$1"
	err := data.QueryRow(queryGet, userID).Scan(&balance)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("There is no user with such ID or wrong query"))
		return
	}

	balanceS := strconv.Itoa(balance)
	w.Write([]byte(balanceS))

	err = data.Close()
	if err != nil {
		log.Println("Cannot close connection")
	}
}

func (h *Handler) AddMoney(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	panic("todo")
}

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	panic("todo")
}

func (h *Handler) SendMoney(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	panic("todo")
}
