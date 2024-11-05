package main

import (
	"blogi/internal/api"
	"blogi/internal/config"
	"blogi/internal/postgres"
	"log"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatalln("Failed to load the config", err.Error())
	}

	db, err := postgres.Connect(conf.DBConn)
	if err != nil {
		log.Fatalln("Failed to connect to the database", err.Error())
	} else if err = db.Migrate(); err != nil {
		log.Fatalln("Failed to do the migrations", err.Error())
	}

	api.Run(api.APIConfig{APIConfig: conf.API, DB: db})
}
