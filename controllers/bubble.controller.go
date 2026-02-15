package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BubbleController struct{}

func (ctrl BubbleController) GetBubbleIndex(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
