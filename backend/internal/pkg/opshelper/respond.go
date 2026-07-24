package opshelper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const internalErrorMessage = "An internal error occurred"

func respondInternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": internalErrorMessage})
}

func respondBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": message})
}

func respondServiceUnavailable(c *gin.Context, message string) {
	c.JSON(http.StatusServiceUnavailable, gin.H{"error": message})
}
