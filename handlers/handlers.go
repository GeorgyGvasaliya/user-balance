package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"user-balance/consts"
	"user-balance/database"
	"user-balance/models"
	"user-balance/utils"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
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
	currency := r.URL.Query().Get("currency")
	db := database.NewPostgresDB(h.cfg)
	defer db.Close()

	query := "select * from public.users where user_id=$1"
	var u models.User
	err := db.QueryRow(query, userID).Scan(&u.UserID, &u.Balance)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(consts.BadUser))
		return
	}

	if currency != "" {
		u.Balance, err = utils.ConvertCurrency(currency, u.Balance)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(consts.BadCurrency))
			return
		}
	}

	if err := json.NewEncoder(w).Encode(u); err != nil {
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
		log.Println(consts.CannotUnmarshal)
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
		log.Println(consts.CannotUnmarshal)
		return
	}

	queryGet := "select * from public.users where user_id=$1"
	var u models.User
	err = db.QueryRow(queryGet, ch.UserID).Scan(&u.UserID, &u.Balance)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(consts.BadUser))
		return
	}

	newBalance := u.Balance - ch.Money
	if newBalance < 0 {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(consts.NoMoney))
		return
	}

	querySet := "UPDATE public.users SET balance=balance-$1 WHERE user_id=$2;"
	db.QueryRow(querySet, ch.Money, ch.UserID)
}

func (h *Handler) SendMoney(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	db := database.NewPostgresDB(h.cfg)
	defer db.Close()
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	var tr models.TransferMoney
	err := json.Unmarshal(body, &tr)
	if err != nil {
		log.Println(consts.CannotUnmarshal)
		return
	}

	// get first user
	queryGet1 := "select * from public.users where user_id=$1"
	var fromUser models.User
	err = db.QueryRow(queryGet1, tr.FromID).Scan(&fromUser.UserID, &fromUser.Balance)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(consts.BadUser))
		return
	}
	// check second user exists
	queryGet2 := "SELECT 1 FROM public.users WHERE user_id=$1"
	var data []byte
	err = db.QueryRow(queryGet2, tr.ToID).Scan(&data)
	if len(data) == 0 {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(consts.BadUser))
		return
	}

	newBalance := fromUser.Balance - tr.Money
	if newBalance < 0 {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(consts.NoMoney))
		return
	}

	tx, err := db.Begin()

	queryWithdraw := "UPDATE public.users SET balance=balance-$1 WHERE user_id=$2;"
	queryAdd := "UPDATE public.users SET balance=balance+$1 WHERE user_id=$2;"

	_, err = tx.Exec(queryWithdraw, tr.Money, tr.FromID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(consts.ServerError))
		tx.Rollback()
		return
	}
	_, err = tx.Exec(queryAdd, tr.Money, tr.ToID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(consts.ServerError))
		tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(consts.ServerError))
		tx.Rollback()
		return
	}
}
