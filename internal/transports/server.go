package transports

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jrodolforojas/libertadfinanciera-backend/internal/services"
)

// WebServer has the logic to start the microservice
type WebServer struct {
}

// StartServer listens and servers this microservice
func (ws *WebServer) StartServer() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	var httpAddr = flag.String("http", fmt.Sprintf(":%s", port), "http listen address")

	flag.Parse()

	ctx := context.Background()

	service := services.NewService()

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		log.Println("listening on port", *httpAddr)
		handler := MakeHTTPHandler(ctx, service)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	log.Println("Server ends ", <-errs)

}
