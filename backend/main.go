package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/database"
	"juraji.nl/chat-quest/providers"
	"juraji.nl/chat-quest/routes"
	"juraji.nl/chat-quest/util"
	"log"
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
		routes.SystemController(apiRouter, db)
		routes.TagsController(apiRouter, db)
		routes.CharactersController(apiRouter, db)
		routes.InstructionPromptsController(apiRouter, db)
		providers.ConnectionProfilesController(apiRouter, db)
		routes.ScenariosController(apiRouter, db)
		routes.ChatPreferencesController(apiRouter, db)
		routes.WorldsController(apiRouter, db)
		routes.ChatSessionsController(apiRouter, db)
		routes.MemoriesController(apiRouter, db)
	}

	// If endpoint is not found, the request is probably a UI resource.
	// Else we just fail ugly, gl hackers.
	log.Printf("Serving Chat Quest UI from directory '%s'", ChatQuestUIRoot)
	router.NoRoute(routes.ChatQuestUIHandler(ChatQuestUIRoot))

	serverAddr := fmt.Sprintf("%s:%s", ApplicationHost, ApplicationPort)
	//goland:noinspection HttpUrlsUsage
	log.Printf("ChatQuest is running on http://%s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatal(err)
	}
}
