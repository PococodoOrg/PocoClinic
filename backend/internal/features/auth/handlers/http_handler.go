package handlers

import (
	"net/http"

	"github.com/dksch/pococlinic/internal/features/auth/commands"
	"github.com/dksch/pococlinic/internal/features/auth/domain"
	"github.com/dksch/pococlinic/internal/features/auth/queries"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles HTTP requests for authentication operations
type AuthHandler struct {
	createUserHandler commands.CreateUserHandler
	loginHandler      commands.LoginHandler
	getUserHandler    queries.GetUserHandler
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(
	createUser commands.CreateUserHandler,
	login commands.LoginHandler,
	getUser queries.GetUserHandler,
) *AuthHandler {
	return &AuthHandler{
		createUserHandler: createUser,
		loginHandler:      login,
		getUserHandler:    getUser,
	}
}

// RegisterRoutes registers the authentication routes with the given router group
func (h *AuthHandler) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.CreateUser)
		auth.POST("/login", h.Login)
		auth.GET("/users/:id", h.GetUser)
	}
}

// CreateUser handles user registration
func (h *AuthHandler) CreateUser(c *gin.Context) {
	var cmd commands.CreateUserCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user, key, err := h.createUserHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		var status int
		var response gin.H

		switch e := err.(type) {
		case *domain.AuthError:
			switch e.Code {
			case domain.ErrEmailTaken:
				status = http.StatusConflict
			default:
				status = http.StatusBadRequest
			}
			response = gin.H{"error": e.Message, "code": e.Code}
		default:
			status = http.StatusInternalServerError
			response = gin.H{"error": "internal server error"}
		}

		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": user,
		"key":  key,
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	var cmd commands.LoginCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Set IP and User-Agent from request
	cmd.IPAddress = c.ClientIP()
	cmd.UserAgent = c.Request.UserAgent()

	session, err := h.loginHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		var status int
		var response gin.H

		switch e := err.(type) {
		case *domain.AuthError:
			switch e.Code {
			case domain.ErrInvalidCredentials:
				status = http.StatusUnauthorized
			case domain.ErrAccountLocked:
				status = http.StatusForbidden
			case domain.ErrUserNotFound:
				status = http.StatusNotFound
			default:
				status = http.StatusBadRequest
			}
			response = gin.H{"error": e.Message, "code": e.Code}
		default:
			status = http.StatusInternalServerError
			response = gin.H{"error": "internal server error"}
		}

		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, session)
}

// GetUser handles user retrieval
func (h *AuthHandler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	query := queries.GetUserQuery{ID: id.String()}
	user, err := h.getUserHandler.Handle(c.Request.Context(), query)
	if err != nil {
		var status int
		var response gin.H

		switch e := err.(type) {
		case *domain.AuthError:
			switch e.Code {
			case domain.ErrUserNotFound:
				status = http.StatusNotFound
			default:
				status = http.StatusBadRequest
			}
			response = gin.H{"error": e.Message, "code": e.Code}
		default:
			status = http.StatusInternalServerError
			response = gin.H{"error": "internal server error"}
		}

		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, user)
}
