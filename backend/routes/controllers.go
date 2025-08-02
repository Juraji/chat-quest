package routes

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func getIDParam(c *gin.Context, key string) (int64, error) {
	idStr := c.Param(key)
	id, err := strconv.ParseInt(idStr, 10, 64)

	return id, err
}

func respondList[T any](c *gin.Context, records []T, err error) {
	if err != nil {
		respondInternalError(c, err)
		log.Print(err)
	} else {
		c.JSON(http.StatusOK, &records)
	}
}

func respondSingle[T any](c *gin.Context, entity *T, err error) {
	if err != nil {
		respondInternalError(c, err)
		log.Print(err)
	} else if entity == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entity not found"})
	} else {
		c.JSON(http.StatusOK, entity)
	}
}

func respondEmpty(c *gin.Context, err error) {
	if err != nil {
		respondInternalError(c, err)
	} else {
		c.Status(http.StatusNoContent)
	}
}

func respondDeleted(c *gin.Context, err error) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete Entity, it might be in use."})
		log.Print(err)
	} else {
		c.JSON(http.StatusOK, gin.H{"error": "Entity was deleted"})
	}
}

func respondInternalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	log.Print(err)
}

func respondBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": message})
}

func respondNotAcceptable(c *gin.Context, message string, err error) {
	c.JSON(http.StatusNotAcceptable, gin.H{"error": message})
	log.Print(err)
}
