package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/characters"
	"juraji.nl/chat-quest/chat-sessions"
	"juraji.nl/chat-quest/cq"
	"juraji.nl/chat-quest/database"
	"juraji.nl/chat-quest/instructions"
	"juraji.nl/chat-quest/memories"
	"juraji.nl/chat-quest/providers"
	"juraji.nl/chat-quest/scenarios"
	"juraji.nl/chat-quest/sse"
	"juraji.nl/chat-quest/system"
	"juraji.nl/chat-quest/util"
	"juraji.nl/chat-quest/worlds"
	"log"
	"os"
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
	rootCtx := context.Background()

	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Failed to close database:", err)
		}
	}(db)

	chatQuestContext := cq.NewChatQuestContext(
		rootCtx,
		db,
		log.New(os.Stdout, "", log.LstdFlags),
	)

	router := gin.New()

	log.Printf("Setting up CORS for hosts: %v", CorsAllowOrigins)
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = CorsAllowOrigins
	router.Use(cors.New(corsConfig))

	if err := router.SetTrustedProxies(GinTrustedProxies); err != nil {
		log.Fatal("Failed to set trusted proxies", err)
	}

	apiRouter := router.Group(ApiBasePath)
	{
		log.Println("Registering routes...")
		system.Routes(chatQuestContext, apiRouter)
		characters.Routes(chatQuestContext, apiRouter)
		instructions.Routes(chatQuestContext, apiRouter)
		providers.Routes(chatQuestContext, apiRouter)
		scenarios.Routes(chatQuestContext, apiRouter)
		worlds.Routes(chatQuestContext, apiRouter)
		chat_sessions.Routes(chatQuestContext, apiRouter)
		memories.Routes(chatQuestContext, apiRouter)
		sse.Routes(apiRouter)
	}

	// If endpoint is not found, the request is probably a UI resource.
	// Else we just fail ugly, gl hackers.
	log.Printf("Serving Chat Quest UI from directory '%s'", ChatQuestUIRoot)
	router.NoRoute(system.ChatQuestUIHandler(ChatQuestUIRoot))

	serverAddr := fmt.Sprintf("%s:%s", ApplicationHost, ApplicationPort)
	//goland:noinspection HttpUrlsUsage
	log.Printf("ChatQuest is running on http://%s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatal(err)
	}
}
