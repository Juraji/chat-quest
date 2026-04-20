package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/api"
	"juraji.nl/chat-quest/core"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
	"juraji.nl/chat-quest/core/ui"
	"juraji.nl/chat-quest/processing"
)

func init() {
	core.InitEnvironment()
	if !core.Env().DebugEnabled {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	core.InitEnvironment()
	env := core.Env()

	// Setup logger
	log.InitLogger(env)

	mainLogger := log.Get()
	mainLogger.Info("System OK! ChatQuest is starting...")

	// Setup DB
	mainLogger.Info("Connecting to database...")
	closeDB := database.InitDB(env)

	defer closeDB()
	mainLogger.Info("Database initialized successfully!")

	// Asynchronous processes (needs to be here, or it won't be compiled!)
	mainLogger.Info("Setting up asynchronous processing...")
	processing.SetupProcessing()

	// New Router!
	router := gin.New()

	// Configure router
	mainLogger.Info("Setting up CORS...", zap.Any("hosts", env.CorsAllowOrigins))
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = env.CorsAllowOrigins
	router.Use(cors.New(corsConfig))

	if err := router.SetTrustedProxies(env.TrustedProxies); err != nil {
		mainLogger.Fatal("Failed to set trusted proxies", zap.Error(err))
	}

	// Register routes (/api)
	apiRouter := router.Group(env.ApiBasePath)
	mainLogger.Info("Registering route handlers...")
	api.SystemRoutes(apiRouter)
	api.PreferencesRoutes(apiRouter)
	api.CharactersRoutes(apiRouter)
	api.InstructionsRoutes(apiRouter)
	api.ProvidersRoutes(apiRouter)
	api.ScenariosRoutes(apiRouter)
	api.SpeciesRoutes(apiRouter)
	api.WorldsRoutes(apiRouter)
	api.ChatSessionsRoutes(apiRouter)
	api.MemoriesRoutes(apiRouter)
	api.SseRoutes(apiRouter)

	// Setup UI host (any non-api route)
	// If endpoint is not found, the request is probably a UI resource, else we just fail ugly.
	mainLogger.Info("Serving ChatQuest web UI", zap.String("root", env.ChatQuestUIRoot))
	router.NoRoute(ui.ChatQuestUIHandler(env.ChatQuestUIRoot))

	// Finally start server
	serverAddr := fmt.Sprintf("%s:%s", env.ApplicationHost, env.ApplicationPort)
	//goland:noinspection HttpUrlsUsage
	mainLogger.Info("ChatQuest is running", zap.String("addr", "http://"+serverAddr))
	if err := router.Run(serverAddr); err != nil {
		mainLogger.Fatal("Failed to start ChatQuest server", zap.Error(err))
	}
}
