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

	var m models.AddMoney
	err := json.Unmarshal(body, &m)
	if err != nil {
		log.Println("Cannot Unmarshal data")
		return
	}

	// insert if not exists, update if exists
	query := "INSERT INTO public.users (user_id, balance) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET balance=users.balance+$2;"
	db.QueryRow(query, m.UserID, m.Money)
}

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//w.Header().Add("Content-Type", "application/json")
	//data := db.NewPostgresDB(h.cfg)
	//body, _ := ioutil.ReadAll(r.Body)
	//
	//defer r.Body.Close()
	//
	//var mr MoneyResponse
	//err := json.Unmarshal(body, &mr)
	//if err != nil {
	//	log.Println("Cannot Unmarshal data")
	//	return
	//}
	//
	//// вообще этот кеш нужен был, чтобы при выводе средств проверять на отрицательный баланс, одним sql запросом не обойдёшься
	//balance := h.cache[mr.Id]
	//newBalance := balance - mr.Money
	//if newBalance < 0 {
	//	w.WriteHeader(http.StatusNotFound)
	//	w.Write([]byte("You have no money to purchase"))
	//	return
	//}
	//
	//h.cache[mr.Id] = newBalance
	//queryUpsert := "INSERT INTO public.wallet (wallet_id, user_id, balance) VALUES ($1, $1, $2) ON CONFLICT (wallet_id) DO UPDATE SET balance=wallet.balance-$2;"
	//data.QueryRow(queryUpsert, mr.Id, mr.Money)
	//
	//err = data.Close()
	//if err != nil {
	//	log.Println("Cannot close connection")
	//}
}

func (h *Handler) SendMoney(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	panic("todo")
}
