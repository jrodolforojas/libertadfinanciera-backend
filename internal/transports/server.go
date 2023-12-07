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
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/configuration"
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
	config, err := configuration.Read()
	if err != nil {
		panic(err)
	}

	logger := utils.NewLogger()
	_ = level.Debug(logger).Log("msg", "service started")

	ctx := context.Background()

	supabaseClient := supabase.InitSupabase(config.Database.SupabaseUrl, config.Database.SupabaseKey)

	_ = level.Debug(logger).Log("msg", "supabase client initialized")

	scrapper := scrapper.NewBCCRScrapper(logger, config.Scrapper)

	_ = level.Debug(logger).Log("msg", "BCCR scrapper initialized")

	repository := supabase.NewSupabase(logger, supabaseClient)

	_ = level.Debug(logger).Log("msg", "repository initialized")

	service := services.NewService(logger, scrapper, repository)

	errs := make(chan error)

	var httpAddr = flag.String("http", fmt.Sprintf(":%s", config.Address.Port), "http listen address")

	flag.Parse()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		_ = level.Debug(logger).Log("msg", "listening", "port", *httpAddr, "transport", "HTTP")
		handler := MakeHTTPHandler(ctx, service, logger)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	_ = level.Debug(logger).Log("shutdown", <-errs)

}
