package main

import (
	"fmt"
	"os"

	"github.com/hugaojanuario/sentinel/internal/router"
)

func main() {
	fmt.Println("Sentinel starting...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	router := router.SetupRouter()
	router.Run(":" + port)
}
