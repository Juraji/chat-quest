package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
)

func getParamAsID(c *gin.Context, key string) (int, bool) {
	idStr := c.Param(key)
	id, err := strconv.ParseInt(idStr, 10, 32)
	return int(id), err == nil
}

func getQueryParamAsIntOr(c *gin.Context, key string, defaultValue int) int {
	str, present := c.GetQuery(key)
	if !present || str == "" {
		return defaultValue
	}

	i, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return defaultValue
	}
	return int(i)
}

func getQueryParamsAsInts(c *gin.Context, key string) ([]int, bool) {
	values, ok := c.GetQueryArray(key)
	if !ok {
		return nil, false
	}

	ints := make([]int, len(values))
	for i, idStr := range values {
		id, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			return nil, false
		}
		ints[i] = int(id)
	}

	return ints, true
}

func respondList[T any](c *gin.Context, list []T, err error) {
	if err != nil {
		respondInternalError(c, err)
	} else {
		c.JSON(http.StatusOK, list)
	}
}

func respondSingle[T any](c *gin.Context, entity *T, err error) {
	if err != nil {
		respondInternalError(c, err)
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

func respondInternalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	if err != nil {
		log.Get().Error("Error during API request", zap.String("uri", c.Request.RequestURI), zap.Error(err))
	} else {
		log.Get().Warn("Not OK API response", zap.String("uri", c.Request.RequestURI))
	}
}

func respondBadRequest(c *gin.Context, message string, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": message})
	log.Get().Warn("Bad API Request",
		zap.String("uri", c.Request.RequestURI),
		zap.String("message", message),
		zap.Error(err))
}

func respondNotAcceptable(c *gin.Context, message string, err error) {
	c.JSON(http.StatusNotAcceptable, gin.H{"error": message})
	log.Get().Warn("Unacceptable data for API request", zap.String("uri", c.Request.RequestURI), zap.Error(err))
}
