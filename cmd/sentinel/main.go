package main

import (
	"fmt"

	"github.com/hugaojanuario/sentinel/internal/api"
)

func main() {
	fmt.Println("Sentinel starting...")
	router := api.NewRouter()

	router.Run(":9191")
}
