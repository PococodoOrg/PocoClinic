package opshelper

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const internalErrorMessage = "An internal error occurred"

func respondInternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": internalErrorMessage})
}

func respondOperationError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": friendlyOperationError(err)})
}

func friendlyOperationError(err error) string {
	if err == nil {
		return internalErrorMessage
	}
	lower := strings.ToLower(err.Error())
	switch {
	case strings.Contains(lower, "database is locked"), strings.Contains(lower, "locked"):
		return "The database is in use. Stop the main PocoClinic app, then try again."
	case strings.Contains(lower, "no such file"), strings.Contains(lower, "not found"):
		return "Backup file not found. Check that it is in your backups folder."
	case strings.Contains(lower, "manifest"), strings.Contains(lower, "checksum"):
		return "This backup file could not be verified. Choose a different file or create a new backup."
	default:
		return "Could not complete the operation. Stop the main app if it is running, then try again."
	}
}

func respondBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": message})
}

func respondServiceUnavailable(c *gin.Context, message string) {
	c.JSON(http.StatusServiceUnavailable, gin.H{"error": message})
}
