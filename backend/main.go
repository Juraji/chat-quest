package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/logging"
	"juraji.nl/chat-quest/core/providers"
	"juraji.nl/chat-quest/core/sse"
	"juraji.nl/chat-quest/core/system"
	"juraji.nl/chat-quest/core/util"
	"juraji.nl/chat-quest/model/characters"
	"juraji.nl/chat-quest/model/chat-sessions"
	"juraji.nl/chat-quest/model/instructions"
	"juraji.nl/chat-quest/model/memories"
	"juraji.nl/chat-quest/model/scenarios"
	"juraji.nl/chat-quest/model/worlds"
)

var (
	ChatQuestUIRoot   = "./browser"
	GinMode           = gin.DebugMode
	GinTrustedProxies []string
	CorsAllowOrigins  = []string{"http://localhost:8080", "http://127.0.0.1:8080"}
	ApplicationHost   = "localhost"
	ApplicationPort   = "8080"
	ApiBasePath       = "/api"
)

func init() {
	util.SetStringFromEnvIfPresent("CHAT_QUEST_UI_ROOT", &ChatQuestUIRoot)
	util.SetStringFromEnvIfPresent("CHAT_QUEST_GIN_MODE", &GinMode)
	util.SetSliceFromEnvIfPresent("CHAT_QUEST_GIN_TRUSTED_PROXIES", &GinTrustedProxies)
	util.SetSliceFromEnvIfPresent("CHAT_QUEST_CORS_ALLOW_ORIGINS", &CorsAllowOrigins)
	util.SetStringFromEnvIfPresent("CHAT_QUEST_APPLICATION_HOST", &ApplicationHost)
	util.SetStringFromEnvIfPresent("CHAT_QUEST_APPLICATION_PORT", &ApplicationPort)
	util.SetStringFromEnvIfPresent("CHAT_QUEST_API_BASE_PATH", &ApiBasePath)

	gin.SetMode(GinMode)
}

func main() {
	logger, err := logging.SetupLogger(false)
	if err != nil {
		panic(err)
	}

	logger.Info("ChatQuest starting...")
	rootCtx := context.Background()

	logger.Info("Connecting to database...")
	db, err := database.InitDB()
	if err != nil {
		logger.Fatal("Failed to initialize database:", zap.Error(err))
	}
	logger.Info("Database initialized successfully!")

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Fatal("Failed to close database:", zap.Error(err))
		}
	}(db)

	chatQuestContext := core.NewChatQuestContext(rootCtx, db, logger)

	router := gin.New()

	logger.Info("Setting up CORS...", zap.Any("hosts", CorsAllowOrigins))
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = CorsAllowOrigins
	router.Use(cors.New(corsConfig))

	if err := router.SetTrustedProxies(GinTrustedProxies); err != nil {
		logger.Fatal("Failed to set trusted proxies", zap.Error(err))
	}

	apiRouter := router.Group(ApiBasePath)
	{
		logger.Info("Registering route handlers...")
		system.Routes(chatQuestContext, apiRouter)
		characters.Routes(chatQuestContext, apiRouter)
		instructions.Routes(chatQuestContext, apiRouter)
		providers.Routes(chatQuestContext, apiRouter)
		scenarios.Routes(chatQuestContext, apiRouter)
		worlds.Routes(chatQuestContext, apiRouter)
		chat_sessions.Routes(chatQuestContext, apiRouter)
		memories.Routes(chatQuestContext, apiRouter)
		sse.Routes(chatQuestContext, apiRouter)
	}

	// If endpoint is not found, the request is probably a UI resource.
	// Else we just fail ugly, gl hackers.
	logger.Info("Serving Chat Quest UI", zap.String("root", ChatQuestUIRoot))
	router.NoRoute(system.ChatQuestUIHandler(chatQuestContext, ChatQuestUIRoot))

	serverAddr := fmt.Sprintf("%s:%s", ApplicationHost, ApplicationPort)
	//goland:noinspection HttpUrlsUsage
	logger.Info("ChatQuest is running", zap.String("addr", "http://"+serverAddr))
	if err := router.Run(serverAddr); err != nil {
		logger.Fatal("Failed to start Chat Quest UI", zap.Error(err))
	}
}
