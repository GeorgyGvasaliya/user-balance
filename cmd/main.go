package main

import (
	"log"
	"net/http"
	"os"
	"user-balance/config"
	"user-balance/database"
	"user-balance/handlers"

	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalln("Config init error", err)
	}

	cfg := database.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	h := handlers.NewHandler(cfg)
	router := httprouter.New()
	h.Register(router)
	log.Fatal(http.ListenAndServe(":"+viper.GetString("server.port"), router))
}
