package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/meteedev/assessment-tax/constant"
	"github.com/meteedev/assessment-tax/tax"
)

func New(){

	service := tax.NewTaxService()

	//add service to handler
	handler := tax.NewHandler(service)
	
	e := echo.New()

	registerRoutes(e,handler)
	
	go startServer(e)
	
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
	fmt.Println(constant.SERVER_GRACEFUL_SHUTDOWN)
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

