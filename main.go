package main

import "github.com/jrodolforojas/libertadfinanciera-backend/internal/transports"

func main() {
	server := transports.WebServer{}
	server.StartServer()
}
