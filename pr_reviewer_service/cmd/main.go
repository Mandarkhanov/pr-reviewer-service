package main

import (
	"database/sql"
	"fmt"
	"pr-reviewer/internal/config"
	"pr-reviewer/internal/handlers"
	"pr-reviewer/internal/repository/postgres"
	"pr-reviewer/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger := initLogger()
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	db, err := connect(cfg.DSN())
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer db.Close()
	logger.Info("Succesfully connected to database",
		zap.String("host", cfg.DBHost),
		zap.String("db", cfg.DBName),
	)

	r := gin.Default()

	repoTeams := postgres.NewTeamRepo()
	repoUsers := postgres.NewUserRepo()
	repoPR := postgres.NewPRRepo()
	svc := service.NewService(db, repoTeams, repoUsers, repoPR)
	handler := handlers.NewHandler(svc)
	handler.InitRoutes(r)

	r.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"db":     "connected",
		})
	})
	logger.Info("Starting server", zap.String("port", cfg.ServerAddress))

	if err := r.Run(cfg.ServerAddress); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}

func initLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	return logger
}

func connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database driver: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
