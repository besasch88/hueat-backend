package main

import (
	"os"
	"time"

	"github.com/hueat/backend/cmd/cli/commands"
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
	"github.com/hueat/backend/internal/pkg/hueat_db"
	"github.com/hueat/backend/internal/pkg/hueat_env"
	"github.com/hueat/backend/internal/pkg/hueat_log"
	"github.com/hueat/backend/internal/pkg/hueat_pubsub"
	"github.com/hueat/backend/internal/pkg/hueat_scheduler"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

/*
This is the entrypoint for the CLI where it is defined the list of
available commands a developer can execute.
Please check in the ´commands´ folder all the available commands.

To execute a command from the main directory of the project
you can run ´go run ./cmd/cli/cli.go <command-name>´
E.g. ´go run ./cmd/cli/cli.go event-replay´
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

	// Init modules
	r := gin.New()
	// Set GIN logger
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	// Init modules
	v1Api := r.Group("cli")
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

	// Create CLI app
	app := cli.NewApp()
	app.Name = "Backend"
	app.Usage = "CLI"

	// Define list of commands available in the CLI
	app.Commands = []cli.Command{
		{
			Name: "event-replay",
			Action: func(c *cli.Context) error {
				return commands.EventReplayCommand(c, pubSubAgent, dbConnection)
			},
			Usage: "Replay historical events optionally filtered by topic and start date",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "start-from",
					Usage:    "Optional ISO 8601 date to start replay from",
					Required: false,
				},
				&cli.StringFlag{
					Name:     "topic-name",
					Usage:    "Optional topic name to filter events",
					Required: false,
				},
			},
		},
		{
			Name: "hash-password",
			Action: func(c *cli.Context) error {
				return commands.HashPasswordCommand(c, dbConnection)
			},
			Usage: "Hash the plain password",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "password",
					Usage:    "The hashed password of the new user",
					Required: true,
				},
			},
		},
	}
	// Start the CLI
	err := app.Run(os.Args)
	if err != nil {
		zap.L().Error("Something went wrong during execution", zap.String("service", "cli"), zap.Error(err))
	}
	// Ensure there is enough time before shutting down the CLI
	// to allow all goroutines to be executed
	zap.L().Info("Shutdown CLI in 3 seconds...", zap.String("service", "webapp"))
	time.Sleep(3 * time.Second)
}
