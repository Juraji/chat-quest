package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/cq"
	"net/http"
	"os"
	"path"
)

func ChatQuestUIHandler(cq *cq.ChatQuestContext, chatQuestUIDir string) func(c *gin.Context) {
	fs := http.FileServer(http.Dir(chatQuestUIDir))

	return func(c *gin.Context) {
		filePath := path.Join(chatQuestUIDir, c.Request.URL.Path)
		_, err := os.Stat(filePath)

		if err == nil {
			fs.ServeHTTP(c.Writer, c.Request)
		} else if os.IsNotExist(err) {
			c.File(path.Join(chatQuestUIDir, "index.html"))
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
			cq.Logger().Error("Failed serving UI file", zap.String("path", filePath))
		}
	}
}
