package routes

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func getID(c *gin.Context, key string) (int32, error) {
	idStr := c.Param(key)
	id, err := strconv.ParseInt(idStr, 10, 32)

	if err != nil {
		return 0, err
	} else {
		return int32(id), nil
	}
}

func respondList[T any](c *gin.Context, records []*T, err error) {
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		log.Fatal(err)
	} else {
		c.JSON(200, &records)
	}
}

func respondSingle[T any](c *gin.Context, entity *T, err error) {
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		log.Fatal(err)
	} else if entity == nil {
		c.JSON(404, gin.H{"error": "Character not found"})
	} else {
		c.JSON(200, entity)
	}
}
