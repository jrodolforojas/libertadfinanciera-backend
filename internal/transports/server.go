package transports

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log/level"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/repositories/supabase"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/services"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/services/scrapper"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/utils"
)

// WebServer has the logic to start the microservice
type WebServer struct {
}

// StartServer listens and servers this microservice
func (ws *WebServer) StartServer() {
	ctx := context.Background()

	logger := utils.NewLogger()
	_ = level.Debug(logger).Log("msg", "service started")

	url := "https://vpnzxyjkngpzghthneea.supabase.co"
	key := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InZwbnp4eWprbmdwemdodGhuZWVhIiwicm9sZSI6ImFub24iLCJpYXQiOjE2OTA2NTMwOTYsImV4cCI6MjAwNjIyOTA5Nn0.AHAU4PPgYMO7FTb3BZCxGwkoZnvawiHgyIODx8W6Seo"
	supabaseClient := supabase.InitSupabase(url, key)

	_ = level.Debug(logger).Log("msg", "supabase client initialized")

	bccr := "https://gee.bccr.fi.cr/indicadoreseconomicos/Cuadros/frmVerCatCuadro.aspx?CodCuadro=400&Idioma=1&FecInicial=%s&FecFinal=%s&Filtro=0"
	scrapper := scrapper.NewBCCRScrapper(logger, bccr)

	_ = level.Debug(logger).Log("msg", "BCCR scrapper initialized")

	repository := supabase.NewSupabase(logger, supabaseClient)

	_ = level.Debug(logger).Log("msg", "repository initialized")

	service := services.NewService(logger, scrapper, repository)

	errs := make(chan error)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	var httpAddr = flag.String("http", fmt.Sprintf(":%s", port), "http listen address")

	flag.Parse()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		_ = level.Debug(logger).Log("msg", "listening", "port", *httpAddr, "transport", "HTTP")
		handler := MakeHTTPHandler(ctx, service)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	_ = level.Debug(logger).Log("shutdown", <-errs)

}
