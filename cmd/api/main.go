package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hugaojanuario/sentinel/internal/http/handler"
	"github.com/hugaojanuario/sentinel/internal/http/router"
	"github.com/hugaojanuario/sentinel/internal/repository"
	"github.com/hugaojanuario/sentinel/internal/services"
	"github.com/hugaojanuario/sentinel/pkg/config"
	"github.com/hugaojanuario/sentinel/pkg/database"
)

func main() {
	fmt.Println("Sentinel starting...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	cfg := config.LoadDotEnv()
	db, err := database.Conn(database.Config{
		DB_HOST:     cfg.DB_HOST,
		DB_PORT:     cfg.DB_PORT,
		DB_USER:     cfg.DB_USER,
		DB_PASSWORD: cfg.DB_PASSWORD,
		DB_NAME:     cfg.DB_NAME,
		DB_SSLMODE:  cfg.DB_SSLMODE,
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	serv := services.NewService(repo)
	auth := handler.NewAuthController(serv)
	router := router.SetupRouter(auth)
	router.Run(":" + port)
}
