package main

import (
	"fmt"

	"github.com/hugaojanuario/sentinel/internal/router"
)

func main() {
	fmt.Println("Sentinel starting...")
	router := router.SetupRouter()
	router.Run(":9090")
}
