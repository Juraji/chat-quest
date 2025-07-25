package main

import (
	"chat-quest/backend/sql"
	sql2 "database/sql"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	db, err := sql.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer func(db *sql2.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Failed to close database:", err)
		}
	}(db)

	router := gin.Default()

	router.GET("/api/characters", func(c *gin.Context) {
		rows, _ := db.Query("SELECT id, name FROM characters")
		defer rows.Close()
		var users []map[string]interface{}
		for rows.Next() {
			var u map[string]interface{}
			// Scan into `u`
			users = append(users, u)
		}
		c.JSON(200, users)
	})

	log.Println("Server running on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
