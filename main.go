package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/wilsonce/connectly-test/api"
	"github.com/wilsonce/connectly-test/dao"
	"github.com/wilsonce/connectly-test/initialize"
	"github.com/wilsonce/connectly-test/tg"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := initialize.InitLogger()
	defer logger.Sync()
	initialize.InitDB()
	//initialize.InitModel()
	//initialize.InitDao()

	dao.SetDefault(initialize.DB)

	botBuilder := tg.NewBotBuilder()
	go botBuilder.AutoBuild()

	r := gin.Default()
	api.InitRoute(r)

	srv := &http.Server{
		Addr:    ":9999",
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			initialize.Logger.Error("HTTP server shutdown with error", zap.Any("error_msg", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	//signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-quit
	initialize.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		initialize.Logger.Error("HTTP server shutdown with error", zap.Any("error_msg", err))
	}
	botBuilder.RemoveAll()
	initialize.Logger.Info("Server gracefully stopped")
}
