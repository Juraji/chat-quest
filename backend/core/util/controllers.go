package util

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	"net/http"
	"strconv"
)

func GetIDParam(c *gin.Context, key string) (int, error) {
	idStr := c.Param(key)
	id, err := strconv.ParseInt(idStr, 10, 64)

	return int(id), err
}

func GetIDsFromQuery(c *gin.Context, key string) ([]int, error) {
	values, ok := c.GetQueryArray(key)
	if !ok {
		return nil, nil
	}

	ids := make([]int, len(values))
	for i, idStr := range values {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, err
		}
		ids[i] = int(id)
	}

	return ids, nil
}

func RespondList[T any](c *gin.Context, records []T, err error) {
	if err != nil {
		RespondInternalError(c, err)
	} else {
		c.JSON(http.StatusOK, &records)
	}
}

func RespondSingle[T any](c *gin.Context, entity *T, err error) {
	if err != nil {
		RespondInternalError(c, err)
	} else if entity == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entity not found"})
	} else {
		c.JSON(http.StatusOK, entity)
	}
}

func RespondEmpty(c *gin.Context, err error) {
	if err != nil {
		RespondInternalError(c, err)
	} else {
		c.Status(http.StatusNoContent)
	}
}

func RespondDeleted(c *gin.Context, err error) {
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Could not delete Entity, it might be in use."})
		log.Get().Error("Conflict during API request, entity could not be deleted", zap.String("uri", c.Request.RequestURI), zap.Error(err))
	} else {
		c.JSON(http.StatusOK, gin.H{"error": "Entity was deleted"})
	}
}

func RespondInternalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	log.Get().Error("Error during API request", zap.String("uri", c.Request.RequestURI), zap.Error(err))
}

func RespondBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": message})
	log.Get().Warn("Bad API Request", zap.String("uri", c.Request.RequestURI), zap.String("message", message))
}

func RespondNotAcceptable(c *gin.Context, message string, err error) {
	c.JSON(http.StatusNotAcceptable, gin.H{"error": message})
	log.Get().Warn("Unacceptable data for API request", zap.String("uri", c.Request.RequestURI), zap.Error(err))
}
