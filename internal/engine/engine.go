package engine

import (
	"fmt"
	"net/http"

	"github.com/alwindoss/ark"
	"github.com/gin-gonic/gin"
)

func Run(cfg *ark.Config) error {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	addr := fmt.Sprintf(":%d", cfg.Port)
	err := r.Run(addr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		return err
	}
	return nil
}
