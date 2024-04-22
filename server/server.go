package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/labstack/echo/v4"
	"github.com/meteedev/assessment-tax/apperrs"
	"github.com/meteedev/assessment-tax/constant"
	"github.com/meteedev/assessment-tax/tax"
)

func New(){

	// Create a logger instance
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	logger = logger.Level(zerolog.InfoLevel)

	// Inject the logger into TaxService
	taxService := tax.NewTaxService(&logger)

	//add service to handler
	handler := tax.NewHandler(taxService)
	
	e := echo.New()

	// config catch error 
	e.Use(apperrs.CustomErrorMiddleware)

	//register rest api route
	registerRoutes(e,handler)
	
	// start servert
	go startServer(e)
	
	//config graceful shutdown
	gracefulShutdownServer(e)
}


func startServer(e *echo.Echo) {
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if err := e.Start(port); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal(err)
	}
}

func gracefulShutdownServer(e *echo.Echo) {
	// Listen for OS signals for graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	//<-ctx.Done()
	fmt.Println(constant.MSG_SERVER_GRACEFUL_SHUTDOWN)
}

// registerRoutes registers all the routes for the application.
func registerRoutes(e *echo.Echo,handler *tax.Handler) {
	
	// Tax routes
	taxGroup := e.Group("/tax")
	taxGroup.POST("/calculations", handler.TaxCalculationsHandler)
	

	// Admin routes
	adminGroup := e.Group("/admin")
	adminGroup.POST("/deductions/personal", handler.DeductionsPersonal)
	

}

