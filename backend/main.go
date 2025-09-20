package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
	"juraji.nl/chat-quest/core/providers"
	"juraji.nl/chat-quest/core/sse"
	"juraji.nl/chat-quest/core/system"
	"juraji.nl/chat-quest/core/util"
	"juraji.nl/chat-quest/model/characters"
	"juraji.nl/chat-quest/model/chat-sessions"
	"juraji.nl/chat-quest/model/instructions"
	"juraji.nl/chat-quest/model/memories"
	"juraji.nl/chat-quest/model/preferences"
	"juraji.nl/chat-quest/model/scenarios"
	"juraji.nl/chat-quest/model/worlds"
	"juraji.nl/chat-quest/processing"
)

var (
	DataDirectory     = "./data"
	ChatQuestUIRoot   = "./browser"
	DebugEnabled      = false
	GinTrustedProxies []string
	CorsAllowOrigins  = []string{"http://localhost:8080", "http://127.0.0.1:8080"}
	ApplicationHost   = "localhost"
	ApplicationPort   = "8080"
	ApiBasePath       = "/api"
)

func init() {
	var err error

	util.SetStringFromEnvIfPresent("CHAT_QUEST_DATA_DIR", &DataDirectory)
	if DataDirectory, err = filepath.Abs(DataDirectory); err != nil {
		panic(errors.Wrapf(err, "Failed to get absolute Path of CHAT_QUEST_DATA_DIR: %s", DataDirectory))
	}
	util.SetStringFromEnvIfPresent("CHAT_QUEST_UI_ROOT", &ChatQuestUIRoot)
	if ChatQuestUIRoot, err = filepath.Abs(ChatQuestUIRoot); err != nil {
		panic(errors.Wrapf(err, "Failed to get absolute Path of CHAT_QUEST_UI_ROOT: %s", ChatQuestUIRoot))
	}

	util.SetSliceFromEnvIfPresent("CHAT_QUEST_GIN_TRUSTED_PROXIES", &GinTrustedProxies)
	util.SetSliceFromEnvIfPresent("CHAT_QUEST_CORS_ALLOW_ORIGINS", &CorsAllowOrigins)
	util.SetStringFromEnvIfPresent("CHAT_QUEST_APPLICATION_HOST", &ApplicationHost)
	util.SetStringFromEnvIfPresent("CHAT_QUEST_APPLICATION_PORT", &ApplicationPort)
	util.SetStringFromEnvIfPresent("CHAT_QUEST_API_BASE_PATH", &ApiBasePath)

	var debugModeVal string
	util.SetStringFromEnvIfPresent("CHAT_QUEST_DEBUG", &debugModeVal)
	if debugModeVal == "true" {
		DebugEnabled = true
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

}

func main() {
	// Setup logger
	log.InitLogger(DebugEnabled, DataDirectory)

	mainLogger := log.Get()
	mainLogger.Info("ChatQuest starting...")

	// Check data dir
	mainLogger.Info("Checking data directory...", zap.String("dataDir", DataDirectory))
	if err := os.MkdirAll(DataDirectory, os.ModePerm); err != nil {
		logger.Error("Failed to create data directory", zap.Error(err))
		return
	}

	// Setup DB
	mainLogger.Info("Connecting to database...")
	closeDB := database.InitDB(DataDirectory)

	defer closeDB()
	mainLogger.Info("Database initialized successfully!")

	// Asynchronous processes (needs to be here, or it won't be compiled!)
	mainLogger.Info("Setting up asynchronous processing...")
	processing.SetupProcessing()

	// New Router!
	router := gin.New()

	// Configure router
	mainLogger.Info("Setting up CORS...", zap.Any("hosts", CorsAllowOrigins))
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = CorsAllowOrigins
	router.Use(cors.New(corsConfig))

	if err := router.SetTrustedProxies(GinTrustedProxies); err != nil {
		mainLogger.Fatal("Failed to set trusted proxies", zap.Error(err))
	}

	// Register routes (/api)
	apiRouter := router.Group(ApiBasePath)
	mainLogger.Info("Registering route handlers...")
	system.Routes(apiRouter)
	preferences.Routes(apiRouter)
	characters.Routes(apiRouter)
	instructions.Routes(apiRouter)
	providers.Routes(apiRouter)
	scenarios.Routes(apiRouter)
	worlds.Routes(apiRouter)
	chat_sessions.Routes(apiRouter)
	memories.Routes(apiRouter)
	sse.Routes(apiRouter)

	// Setup UI host (any non-api route)
	// If endpoint is not found, the request is probably a UI resource, else we just fail ugly.
	mainLogger.Info("Serving ChatQuest web UI", zap.String("root", ChatQuestUIRoot))
	router.NoRoute(system.ChatQuestUIHandler(ChatQuestUIRoot))

	// Finally start server
	serverAddr := fmt.Sprintf("%s:%s", ApplicationHost, ApplicationPort)
	//goland:noinspection HttpUrlsUsage
	mainLogger.Info("ChatQuest is running", zap.String("addr", "http://"+serverAddr))
	if err := router.Run(serverAddr); err != nil {
		mainLogger.Fatal("Failed to start ChatQuest server", zap.Error(err))
	}
}
