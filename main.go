package main

import (
	"actlabs-auth/handler"
	"actlabs-auth/middleware"
	"actlabs-auth/repository"
	"actlabs-auth/service"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

func main() {

	// Get the log level from the environment or default to 8.
	// logLevel := os.Getenv("LOG_LEVEL")
	// logLevelInt, err := strconv.Atoi(logLevel)
	// if err != nil {
	// 	logLevelInt = 8
	// }

	// Create a new logger.
	opts := slog.HandlerOptions{
		//Level:     slog.Level(logLevelInt),
		AddSource: true,
	}

	slog.SetDefault(slog.New(opts.NewJSONHandler(os.Stderr)))

	router := gin.Default()
	router.SetTrustedProxies(nil)

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	router.Use(cors.New(config))

	authRouter := router.Group("/")

	authRouter.Use(middleware.AuthRequired())

	authService := service.NewAuthService(repository.NewAuthRepository())
	handler.NewAuthHandler(authRouter, authService)

	router.Run()
}
