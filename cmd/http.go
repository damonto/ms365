package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/damonto/msonline/internal/app"
	"github.com/damonto/msonline/internal/pkg/config"
	"github.com/damonto/msonline/internal/pkg/logger"
	"github.com/spf13/cobra"
)

var httpCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the microsoft online RESTful API",
	Run: func(cmd *cobra.Command, args []string) {
		srv := &http.Server{
			Addr:         config.Cfg.App.ListenAddr,
			Handler:      app.Handler(),
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		}

		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Sugar.Fatalf("Listen error: %v", err)
			}
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-quit

		logger.Sugar.Info("Shutdown server...")
		var timeout time.Duration
		if os.Getenv("GIN_MODE") != "release" {
			timeout = 0
		} else {
			timeout = 10 * time.Second
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Sugar.Fatalf("Server shutdown error %v", err)
		}

		select {
		case <-ctx.Done():
			logger.Sugar.Info("Server exiting")
		}
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}
