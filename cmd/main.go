package main

import (
	"fmt"
	"log"
	"os"

	"github.com/iki-rumondor/go-speech/internal/config"
	"github.com/iki-rumondor/go-speech/internal/migrate"
	"github.com/iki-rumondor/go-speech/internal/routes"
)

func main() {
	gormDB, err := config.NewMysqlDB()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	if len(os.Args)-1 > 0 {
		migrate.ReadTerminal(gormDB, os.Args)
		return
	}

	handlers := config.GetAppHandlers(gormDB)

	var PORT = os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	routes.StartServer(handlers).Run(fmt.Sprintf(":%s", PORT))
}
