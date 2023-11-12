package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/boskuv/finance-control_expenses-service/internal/expenses"
	dbx "github.com/go-ozzo/ozzo-dbx"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/cors"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	//"github.com/qiangxue/go-rest-api/internal/auth"

	//"github.com/qiangxue/go-rest-api/internal/healthcheck"

	//"github.com/qiangxue/go-rest-api/pkg/accesslog"
	"github.com/boskuv/finance-control_expenses-service/pkg/dbcontext"
	"github.com/boskuv/finance-control_expenses-service/pkg/log"
)

// Version indicates the current version of the application.
var Version = "1.0.0"

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

func main() {
	flag.Parse()
	// create root logger tagged with server version
	logger := log.New().With(nil, "version", Version)

	// load application configurations
	// cfg := config.MustLoad()
	// cfg, err := config.Load(*flagConfig, logger)
	// if err != nil {
	// 	logger.Errorf("failed to load application configuration: %s", err)
	// 	os.Exit(-1)
	// }

	// connect to the database
	db, err := dbx.MustOpen("sqlite3", ":memory:") //cfg.DSN)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
	db.QueryLogFunc = logDBQuery(logger)
	db.ExecLogFunc = logDBExec(logger)
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(err)
		}
	}()

	// build HTTP server
	address := fmt.Sprintf(":%v", "8080") //cfg.ServerPort)
	hs := &http.Server{
		Addr:    address,
		Handler: buildHandler(logger, dbcontext.New(db)), // cfg
	}

	// start the HTTP server with graceful shutdown
	go routing.GracefulShutdown(hs, 10*time.Second, logger.Infof)
	logger.Infof("server %v is running at %v", Version, address)
	if err := hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(err)
		os.Exit(-1)
	}
}

// buildHandler sets up the HTTP routing and builds an HTTP handler.
func buildHandler(logger log.Logger, db *dbcontext.DB) http.Handler { // , cfg *config.Config
	router := routing.New()

	router.Use(
		//accesslog.Handler(logger),
		//errors.Handler(logger),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.AllowAll),
	)

	//healthcheck.RegisterHandlers(router, Version)

	rg := router.Group("/v1")

	//authHandler := auth.Handler(cfg.JWTSigningKey)

	expenses.RegisterHandlers(rg.Group(""),
		expenses.NewService(expenses.NewRepository(db, logger), logger),
		logger,
	)

	// auth.RegisterHandlers(rg.Group(""),
	// 	auth.NewService(cfg.JWTSigningKey, cfg.JWTExpiration, logger),
	// 	logger,
	// )

	return router
}

// logDBQuery returns a logging function that can be used to log SQL queries.
func logDBQuery(logger log.Logger) dbx.QueryLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB query successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB query error: %v", err)
		}
	}
}

// logDBExec returns a logging function that can be used to log SQL executions.
func logDBExec(logger log.Logger) dbx.ExecLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB execution successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB execution error: %v", err)
		}
	}
}

// func init() {
// 	viper.AddConfigPath("config")
// 	viper.SetConfigName("config")
// 	viper.SetConfigType("yaml")

// 	if err := viper.ReadInConfig(); err != nil {
// 		log.Fatal("Error reading config file", err)
// 	}
// }

// const (
// 	envLocal = "local"
// 	envDev   = "dev"
// 	envProd  = "prod"
// )

// func main() {
// 	cfg := config.MustLoad()

// 	log := setupLogger(cfg.Env)
// 	log = log.With(slog.String("env", cfg.Env))

// 	log.Info("initializing server", slog.String("address", "Not implemented"))
// 	log.Debug("logger debug mode enabled")
// }

// func setupLogger(env string) *slog.Logger {
// 	var log *slog.Logger

// 	switch env {
// 	case envLocal:
// 		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
// 	case envDev:
// 		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
// 	case envProd:
// 		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
// 	}

// 	return log
// }
