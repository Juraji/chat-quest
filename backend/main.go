package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/model"
	"juraji.nl/chat-quest/routes"
	"log"
	"os"
)

var (
	ChatQuestUIDir = ""
	GinMode        = gin.DebugMode
)

func init() {
	gin.SetMode(GinMode)

	if ChatQuestUIDir == "" {
		ChatQuestUIDir = os.Getenv("CHAT_QUEST_UI_DIR")
	}
}

func main() {
	db, err := model.InitDB()
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

	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatal("Failed to set trusted proxies", err)
	}

	// Register API routes first
	apiRouter := router.Group("/api")
	{
		log.Println("Registering routes...")
		routes.TagsController(apiRouter, db)
		routes.CharactersController(apiRouter, db)
	}

	// Add a custom handler for static files that don't match our API pattern
	log.Printf("Serving Chat Quest UI from directory '%s'", ChatQuestUIDir)
	router.NoRoute(routes.ChatQuestUIHandler(ChatQuestUIDir))

	log.Println("ChatQuest is running on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
