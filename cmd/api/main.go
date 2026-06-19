package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hugaojanuario/sentinel/internal/alerting"
	"github.com/hugaojanuario/sentinel/internal/healthcheck"
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
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	checker := healthcheck.NewChecker(cfg.CHECK_INTERVAL)
	alerter := alerting.NewAlerter(cfg.TELEGRAM_BOT_TOKEN, cfg.TELEGRAM_CHAT_ID, checker.Events())

	go alerter.Run(ctx)
	go checker.Run(ctx)

	secret := []byte(cfg.JWT_SECRET)
	repo := repository.NewRepository(db)
	serv := services.NewService(repo, secret)
	auth := handler.NewAuthController(serv)
	router := router.SetupRouter(auth, secret)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("erro no servidor: %v", err)
		}
	}()
	log.Printf("Sentinel ouvindo na porta %s", port)

	<-ctx.Done()
	log.Println("desligando servidor...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("erro no shutdown: %v", err)
	}

	log.Println("servidor encerrado")
}
