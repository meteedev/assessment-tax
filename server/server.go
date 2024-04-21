package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/meteedev/assessment-tax/constant"
)

func New(){
	e := echo.New()

	// e.GET("/api/v1/wallets", handler.WalletHandler)

	//graceful shut down
	go func() {
		e.Logger.Debug("in check graceful")
		port := fmt.Sprintf(":%s",os.Getenv("PORT"))
		if err := e.Start(port); err != nil && err != http.ErrServerClosed { // Start server
			e.Logger.Fatal(err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	<-ctx.Done()
	fmt.Println(constant.SERVER_GRACEFUL_SHUTDOWN)

}


