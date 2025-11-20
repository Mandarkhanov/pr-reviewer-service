package main

import (
	"pr-reviewer/pkg/config"
	"pr-reviewer/pkg/database"

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

	db, err := database.Connect(cfg.DSN())
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer db.Close()
	logger.Info("Succesfully connected to database",
		zap.String("host", cfg.DBHost),
		zap.String("db", cfg.DBName),
	)

	r := gin.Default()
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
