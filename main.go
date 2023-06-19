package main

import (
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/services"
)

func main() {
	service := services.NewService()
	service.GetDolarColonesChange()
}
