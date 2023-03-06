package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Controller) Healthz(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	c.String(http.StatusOK, "OK")
}
