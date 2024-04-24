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
	"github.com/labstack/echo/v4/middleware"
	"github.com/meteedev/assessment-tax/apperrs"
	"github.com/meteedev/assessment-tax/authen"
	"github.com/meteedev/assessment-tax/constant"
	"github.com/meteedev/assessment-tax/postgres"
	"github.com/meteedev/assessment-tax/tax/handler"
	"github.com/meteedev/assessment-tax/tax/repository"
	"github.com/meteedev/assessment-tax/tax/service"
	"github.com/rs/zerolog"
)

func New(){

	db , err := postgres.NewDb()
	if err != nil {
	 	panic(err)
	}
	


	logger := zerolog.New(os.Stdout).With().Logger()
	logger = logger.Level(zerolog.DebugLevel)

	// inject db to repository
	taxDeductConfigRepo := repository.NewTaxDeductConfigRepo(db)

	// Inject the logger into TaxService
	taxService := service.NewTaxService(&logger,taxDeductConfigRepo)

	//add service to handler
	handler := handler.NewTaxHandler(taxService)
	
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
func registerRoutes(e *echo.Echo,handler *handler.TaxHandler) {
	
	// Tax routes
	taxGroup := e.Group("/tax")
	taxGroup.POST("/calculations", handler.TaxCalculationsHandler)
	

	// Admin routes
	adminGroup := e.Group("/admin")
	adminGroup.Use(middleware.BasicAuth(authen.AuthMiddleware))
	adminGroup.POST("/deductions/personal", handler.DeductionsPersonal)
	adminGroup.POST("/deductions/k-receipt", handler.DeductionsKreceipt)

}

