package main

import (
	"chat-quest/backend/model"
	"chat-quest/backend/routes"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
)

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

	// Set only to trust localhost proxies
	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatal("Failed to set trusted proxies", err)
	}

	// Create a new router group with the /api prefix
	apiRouter := router.Group("/api")
	{
		log.Println("Registering routes...")
		routes.CharactersController(apiRouter, db)
	}

	log.Println("Server running on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
