package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hueat/backend/internal/app/auth"
	"github.com/hueat/backend/internal/app/healthCheck"
	"github.com/hueat/backend/internal/app/menu"
	"github.com/hueat/backend/internal/app/menuCategory"
	"github.com/hueat/backend/internal/app/menuItem"
	"github.com/hueat/backend/internal/app/menuOption"
	"github.com/hueat/backend/internal/app/order"
	"github.com/hueat/backend/internal/app/printer"
	"github.com/hueat/backend/internal/app/statistics"
	"github.com/hueat/backend/internal/app/table"
	"github.com/hueat/backend/internal/pkg/hueat_auth"
	"github.com/hueat/backend/internal/pkg/hueat_cors"
	"github.com/hueat/backend/internal/pkg/hueat_db"
	"github.com/hueat/backend/internal/pkg/hueat_env"
	"github.com/hueat/backend/internal/pkg/hueat_log"
	"github.com/hueat/backend/internal/pkg/hueat_pubsub"
	"github.com/hueat/backend/internal/pkg/hueat_router"
	"github.com/hueat/backend/internal/pkg/hueat_scheduler"
	ginzap "github.com/gin-contrib/zap"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

/*
This is the entrypoint for the Webapp application built on top of
the GIN framework. It exposes a set of APIs.

To start it you can run ´go run ./cmd/webapp/main.go´
*/
func main() {
	// Set default Timezone
	os.Setenv("TZ", "UTC")
	// ENV Variables
	envs := hueat_env.ReadEnvs()
	// Set Logger
	logger := hueat_log.NewLogger(envs.AppMode)
	zap.ReplaceGlobals(logger)
	// DB Connection
	dbConnection := hueat_db.NewDatabaseConnection(
		envs.DbHost,
		envs.DbUsername,
		envs.DbPassword,
		envs.DbName,
		envs.DbPort,
		envs.DbSslMode,
		envs.DbLogSlowQueryThreshold,
		envs.AppMode,
	)
	// Scheduler
	scheduler := hueat_scheduler.NewScheduler()
	// PUB-SUB agent
	pubSubAgent := hueat_pubsub.NewPubSubAgent(dbConnection, scheduler, envs.PubSubPersistEventsOnDb, envs.PubSubPersistEventsRetentionDays, envs.PubSubSyncMode)

	// Start Server
	zap.L().Info("Starting HTTP Server...", zap.String("service", "webapp"))
	gin.SetMode(envs.AppMode)
	r := gin.New()
	r.SetTrustedProxies(nil)
	// Set GIN logger
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	// Cors Middleware
	allowOrigins := []string{envs.AppCorsOrigin}
	if envs.AppMode != "release" {
		allowOrigins = append(allowOrigins, hueat_cors.LocalhostOrigin)
	}
	r.Use(hueat_cors.CorsMiddleware(allowOrigins))

	// Init Authentication middleware
	authConfig := hueat_auth.AuthConfig{
		JwtSecret: envs.AuthJwtSecret,
	}
	hueat_auth.InitAuthMiddleware(authConfig)

	r.NoRoute(func(ctx *gin.Context) {
		hueat_router.ReturnNotFoundError(ctx, errors.New("endpoint-not-found"))
	})

	// Init moduels that will start exposing endpoints and consumers of internal events
	v1Api := r.Group("api/v1")
	healthCheck.Init(envs, dbConnection, v1Api)
	auth.Init(envs, dbConnection, scheduler, v1Api)
	printer.Init(envs, dbConnection, pubSubAgent, v1Api)
	menuCategory.Init(envs, dbConnection, pubSubAgent, v1Api)
	menuItem.Init(envs, dbConnection, pubSubAgent, v1Api)
	menuOption.Init(envs, dbConnection, pubSubAgent, v1Api)
	menu.Init(envs, dbConnection, pubSubAgent, v1Api)
	table.Init(envs, dbConnection, pubSubAgent, v1Api)
	order.Init(envs, dbConnection, pubSubAgent, v1Api)
	statistics.Init(envs, dbConnection, pubSubAgent, v1Api)

	// Start the scheduler
	if err := scheduler.Init(); err != nil {
		panic(err)
	}

	// Start the application
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", envs.AppPort),
		Handler: r,
	}

	go func() {
		// Start the HTTP Server and listen for errors
		zap.L().Info(fmt.Sprintf("HTTP Server started on port %d", envs.AppPort), zap.String("service", "webapp"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Error("Server Startup Error", zap.String("service", "webapp"), zap.Error(err))
			panic(err)
		}
	}()

	/*
		Wait for interrupt Signals to gracefully shutdown the server
		with a timeout of 3 seconds to ensure all the connection are closed
		and all the pubsub chain activities are performed without receiving
		any additional http request
	*/
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	zap.L().Info("Shutdown Server in 3 seconds...", zap.String("service", "webapp"))

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	hueat_db.CloseDatabaseConnection(dbConnection)
	scheduler.Close()
	pubSubAgent.Close()
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("Server Shutdown Error", zap.String("service", "webapp"), zap.Error(err))
	}

	<-ctx.Done()
	zap.L().Info("Server exited!", zap.String("service", "webapp"))
}
