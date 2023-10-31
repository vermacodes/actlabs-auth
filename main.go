package main

import (
	"actlabs-auth/handler"
	"actlabs-auth/middleware"
	"actlabs-auth/repository"
	"actlabs-auth/service"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

func main() {

	// Get the log level from the environment or default to 8.
	logLevel := os.Getenv("LOG_LEVEL")
	logLevelInt, err := strconv.Atoi(logLevel)
	if err != nil {
		logLevelInt = 0
	}

	// Create a new logger.
	opts := slog.HandlerOptions{
		Level:     slog.Level(logLevelInt),
		AddSource: true,
	}

	slog.SetDefault(slog.New(opts.NewTextHandler(os.Stderr)))

	router := gin.Default()
	router.SetTrustedProxies(nil)

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173", "https://ashisverma.z13.web.core.windows.net", "https://actlabs.z13.web.core.windows.net", "https://actlabsbeta.z13.web.core.windows.net", "https://actlabs.azureedge.net", "https://*.azurewebsites.net"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}

	router.Use(cors.New(config))

	// auth required
	authRouter := router.Group("/")

	authRouter.Use(middleware.AuthRequired())

	labService := service.NewLabService(repository.NewLabRepository())
	handler.NewLabHandler(authRouter, labService)

	authService := service.NewAuthService(repository.NewAuthRepository())
	handler.NewAuthHandler(authRouter, authService)

	loggingService := service.NewLoggingService(repository.NewLoggingRepository())
	handler.NewLoggingHandler(authRouter, loggingService)

	adminAuthRouter := authRouter.Group("/")
	adminAuthRouter.Use(middleware.AdminRequired(authService))
	handler.NewAdminAuthHandler(adminAuthRouter, authService)

	// mentor required
	mentorAuthRouter := authRouter.Group("/")
	mentorAuthRouter.Use(middleware.MentorRequired(authService))
	handler.NewLabHandlerMentorRequired(mentorAuthRouter, labService)

	// apply middleware to all POST requests.
	labMentorAuthRouter := mentorAuthRouter.Group("/")
	labMentorAuthRouter.Use(middleware.UpdateCredits())
	handler.NewLabHandlerMentorRequiredWithCreditMiddleware(labMentorAuthRouter, labService)

	assignmentService := service.NewAssignmentService(repository.NewAssignmentRepository(), labService)
	handler.NewAssignmentHandler(authRouter, assignmentService)

	router.Run()
}
