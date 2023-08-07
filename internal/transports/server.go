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

	"github.com/jrodolforojas/libertadfinanciera-backend/internal/repositories/supabase"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/services"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/services/scrapper"
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

	url := "https://vpnzxyjkngpzghthneea.supabase.co"
	key := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InZwbnp4eWprbmdwemdodGhuZWVhIiwicm9sZSI6ImFub24iLCJpYXQiOjE2OTA2NTMwOTYsImV4cCI6MjAwNjIyOTA5Nn0.AHAU4PPgYMO7FTb3BZCxGwkoZnvawiHgyIODx8W6Seo"
	supabaseClient := supabase.InitSupabase(url, key)

	scrapper := scrapper.NewBCCRScrapper("https://gee.bccr.fi.cr/indicadoreseconomicos/Cuadros/frmVerCatCuadro.aspx?CodCuadro=400&Idioma=1&FecInicial=%s&FecFinal=%s&Filtro=0")
	repository := supabase.NewSupabase(supabaseClient)

	service := services.NewService(scrapper, repository)

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
