package handlers

import (
	"net/http"

	auditdomain "github.com/dksch/pococlinic/internal/features/audit/domain"
	auditqueries "github.com/dksch/pococlinic/internal/features/audit/queries"
	"github.com/dksch/pococlinic/internal/features/auth/commands"
	"github.com/dksch/pococlinic/internal/features/auth/domain"
	authinfra "github.com/dksch/pococlinic/internal/features/auth/infrastructure"
	"github.com/dksch/pococlinic/internal/features/auth/middleware"
	"github.com/dksch/pococlinic/internal/features/auth/queries"
	"github.com/dksch/pococlinic/internal/pkg/config"
	"github.com/dksch/pococlinic/internal/pkg/httperr"
	"github.com/dksch/pococlinic/internal/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles HTTP requests for authentication operations
type AuthHandler struct {
	logger              *logging.Logger
	authConfig          config.AuthConfig
	createUserHandler   commands.CreateUserHandler
	loginHandler        commands.LoginHandler
	refreshTokenHandler commands.RefreshTokenHandler
	logoutHandler       commands.LogoutHandler
	reissueBadgeHandler commands.ReissueBadgeHandler
	changePINHandler    commands.ChangePINHandler
	unlockUserHandler   commands.UnlockUserHandler
	updateUserHandler   commands.UpdateUserHandler
	deleteUserHandler   commands.DeleteUserHandler
	getUserHandler      queries.GetUserHandler
	listUsersHandler    queries.GetUsersHandler
	getUserAuditHandler auditqueries.GetUserAuditHandler
	auditLogger         auditdomain.Logger
	bootstrapDataDir    string
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(
	logger *logging.Logger,
	authConfig config.AuthConfig,
	bootstrapDataDir string,
	createUser commands.CreateUserHandler,
	login commands.LoginHandler,
	refresh commands.RefreshTokenHandler,
	logout commands.LogoutHandler,
	reissueBadge commands.ReissueBadgeHandler,
	changePIN commands.ChangePINHandler,
	unlockUser commands.UnlockUserHandler,
	updateUser commands.UpdateUserHandler,
	deleteUser commands.DeleteUserHandler,
	getUser queries.GetUserHandler,
	listUsers queries.GetUsersHandler,
	getUserAudit auditqueries.GetUserAuditHandler,
	auditLogger auditdomain.Logger,
) *AuthHandler {
	return &AuthHandler{
		logger:              logger,
		authConfig:          authConfig,
		createUserHandler:   createUser,
		loginHandler:        login,
		refreshTokenHandler: refresh,
		logoutHandler:       logout,
		reissueBadgeHandler: reissueBadge,
		changePINHandler:    changePIN,
		unlockUserHandler:   unlockUser,
		updateUserHandler:   updateUser,
		deleteUserHandler:   deleteUser,
		getUserHandler:      getUser,
		listUsersHandler:    listUsers,
		getUserAuditHandler: getUserAudit,
		auditLogger:         auditLogger,
		bootstrapDataDir:    bootstrapDataDir,
	}
}

// RegisterRoutes registers the authentication routes with the given router group
func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware, loginLimiter gin.HandlerFunc) {
	auth := router.Group("/auth")
	{
		auth.POST("/login/staff", loginLimiter, h.StaffLogin)
		auth.POST("/login/admin", loginLimiter, h.AdminLogin)
		auth.POST("/refresh", loginLimiter, h.Refresh)
		auth.POST("/logout", h.Logout)
		auth.GET("/me", authMiddleware.RequireAuth(), h.Me)
		auth.POST("/change-pin", authMiddleware.RequireAuth(), h.ChangePIN)
		auth.POST("/register", authMiddleware.RequireAuth(), authMiddleware.RequirePINChanged(), authMiddleware.RequireRole(domain.RoleAdmin), h.CreateUser)
		auth.GET("/users", authMiddleware.RequireAuth(), authMiddleware.RequirePINChanged(), authMiddleware.RequireRole(domain.RoleAdmin), h.ListUsers)
		auth.GET("/users/:id", authMiddleware.RequireAuth(), authMiddleware.RequirePINChanged(), authMiddleware.RequireRole(domain.RoleAdmin), h.GetUser)
		auth.GET("/users/:id/audit-logs", authMiddleware.RequireAuth(), authMiddleware.RequirePINChanged(), authMiddleware.RequireRole(domain.RoleAdmin), h.GetUserAuditLogs)
		auth.PUT("/users/:id", authMiddleware.RequireAuth(), authMiddleware.RequirePINChanged(), authMiddleware.RequireRole(domain.RoleAdmin), h.UpdateUser)
		auth.POST("/users/:id/unlock", authMiddleware.RequireAuth(), authMiddleware.RequirePINChanged(), authMiddleware.RequireRole(domain.RoleAdmin), h.UnlockUser)
		auth.DELETE("/users/:id", authMiddleware.RequireAuth(), authMiddleware.RequirePINChanged(), authMiddleware.RequireRole(domain.RoleAdmin), h.DeleteUser)
		auth.POST("/users/:id/reissue-badge", authMiddleware.RequireAuth(), authMiddleware.RequirePINChanged(), authMiddleware.RequireRole(domain.RoleAdmin), h.ReissueBadge)
	}
}

// CreateUser handles user registration
func (h *AuthHandler) CreateUser(c *gin.Context) {
	var cmd commands.CreateUserCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}

	user, key, err := h.createUserHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		httperr.RespondUserManagement(c, h.logger, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": toStaffUserJSON(user),
		"key":  key,
	})

	h.auditUserEvent(c, auditdomain.EventUserCreated, user.ID.String(), true, map[string]string{
		"email": user.Email,
		"role":  string(user.Role),
	})
}

// StaffLogin handles badge + PIN authentication for daily staff use
func (h *AuthHandler) StaffLogin(c *gin.Context) {
	var cmd commands.StaffLoginCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}

	cmd.IPAddress = c.ClientIP()
	cmd.UserAgent = c.Request.UserAgent()

	response, err := h.loginHandler.HandleStaff(c.Request.Context(), cmd)
	if err != nil {
		httperr.RespondAuth(c, h.logger, err)
		return
	}

	setAuthCookies(c, h.authConfig, response.AccessToken, response.RefreshToken)
	c.JSON(http.StatusOK, gin.H{
		"user": response.User,
	})
}

// AdminLogin handles email + key + PIN authentication for bootstrap and administrator setup
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var cmd commands.AdminLoginCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}

	cmd.IPAddress = c.ClientIP()
	cmd.UserAgent = c.Request.UserAgent()

	response, err := h.loginHandler.HandleAdmin(c.Request.Context(), cmd)
	if err != nil {
		httperr.RespondAuth(c, h.logger, err)
		return
	}

	setAuthCookies(c, h.authConfig, response.AccessToken, response.RefreshToken)
	if h.bootstrapDataDir != "" {
		_ = authinfra.RemoveBootstrapCredentialsFile(h.bootstrapDataDir)
	}
	c.JSON(http.StatusOK, gin.H{
		"user": response.User,
	})
}

// Refresh issues a new access token using the refresh token cookie
func (h *AuthHandler) Refresh(c *gin.Context) {
	cmd := commands.RefreshTokenCommand{
		RefreshToken: readRefreshTokenCookie(c, h.authConfig),
		IPAddress:    c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
	}

	response, err := h.refreshTokenHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		clearAuthCookies(c, h.authConfig)
		httperr.RespondAuth(c, h.logger, err)
		return
	}

	setAuthCookies(c, h.authConfig, response.AccessToken, response.RefreshToken)
	c.JSON(http.StatusOK, gin.H{
		"user": response.User,
	})
}

// Logout ends the current session
func (h *AuthHandler) Logout(c *gin.Context) {
	cmd := commands.LogoutCommand{
		RefreshToken: readRefreshTokenCookie(c, h.authConfig),
		IPAddress:    c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
	}

	_ = h.logoutHandler.Handle(c.Request.Context(), cmd)
	clearAuthCookies(c, h.authConfig)
	c.Status(http.StatusNoContent)
}

// Me returns the currently authenticated user
func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		httperr.Respond(c, h.logger, domain.NewAuthError(domain.ErrInvalidToken, "missing user context"))
		return
	}

	query := queries.GetUserQuery{ID: userID.(string)}
	user, err := h.getUserHandler.Handle(c.Request.Context(), query)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, toStaffUserJSON(user))
}

// ChangePIN updates the authenticated user's PIN.
func (h *AuthHandler) ChangePIN(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		httperr.Respond(c, h.logger, domain.NewAuthError(domain.ErrInvalidToken, "missing user context"))
		return
	}

	var cmd commands.ChangePINCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}

	cmd.UserID = userID.(string)
	if sessionID, ok := c.Get("sessionID"); ok {
		cmd.SessionID = sessionID.(string)
	}
	cmd.IPAddress = c.ClientIP()
	cmd.UserAgent = c.Request.UserAgent()

	user, err := h.changePINHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		httperr.RespondAuth(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, toStaffUserJSON(user))
}

// ReissueBadge rotates a user's badge key for administrators.
func (h *AuthHandler) ReissueBadge(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid user id")
		return
	}

	actorID, _ := c.Get("userID")

	cmd := commands.ReissueBadgeCommand{
		UserID:    id.String(),
		ActorID:   actorID.(string),
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}

	user, key, err := h.reissueBadgeHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		httperr.RespondUserManagement(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": toStaffUserJSON(user),
		"key":  key,
	})
}

// UnlockUser clears a lockout after failed sign-in attempts.
func (h *AuthHandler) UnlockUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid user id")
		return
	}

	actorID, _ := c.Get("userID")
	cmd := commands.UnlockUserCommand{
		UserID:    id.String(),
		ActorID:   actorID.(string),
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}

	user, err := h.unlockUserHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		httperr.RespondUserManagement(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, toStaffUserJSON(user))
}

// ListUsers handles retrieving a paginated list of users
func (h *AuthHandler) ListUsers(c *gin.Context) {
	var query queries.GetUsersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid query parameters")
		return
	}

	result, err := h.listUsersHandler.Handle(c.Request.Context(), query)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}

	users := make([]StaffUserJSON, len(result.Users))
	for i, user := range result.Users {
		users[i] = toStaffUserJSON(user)
	}

	c.JSON(http.StatusOK, gin.H{
		"users":       users,
		"totalCount":  result.TotalCount,
		"currentPage": result.CurrentPage,
		"pageSize":    result.PageSize,
		"totalPages":  result.TotalPages,
	})
}

// GetUser handles user retrieval
func (h *AuthHandler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid user id")
		return
	}

	query := queries.GetUserQuery{ID: id.String()}
	user, err := h.getUserHandler.Handle(c.Request.Context(), query)
	if err != nil {
		httperr.RespondUserManagement(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, toStaffUserJSON(user))
}

// UpdateUser handles staff profile updates.
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid user id")
		return
	}

	var cmd commands.UpdateUserCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid request body")
		return
	}

	actorID, _ := c.Get("userID")
	cmd.UserID = id.String()
	cmd.ActorID = actorID.(string)
	cmd.IPAddress = c.ClientIP()
	cmd.UserAgent = c.Request.UserAgent()

	user, err := h.updateUserHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		httperr.RespondUserManagement(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, toStaffUserJSON(user))
}

// DeleteUser removes a staff account.
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid user id")
		return
	}

	actorID, _ := c.Get("userID")
	cmd := commands.DeleteUserCommand{
		UserID:    id.String(),
		ActorID:   actorID.(string),
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}

	if err := h.deleteUserHandler.Handle(c.Request.Context(), cmd); err != nil {
		httperr.RespondUserManagement(c, h.logger, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetUserAuditLogs returns login and action history for a staff member.
func (h *AuthHandler) GetUserAuditLogs(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httperr.RespondValidation(c, h.logger, "invalid user id")
		return
	}

	var query auditqueries.GetUserAuditQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httperr.RespondValidation(c, h.logger, "invalid query parameters")
		return
	}
	query.UserID = id.String()

	result, err := h.getUserAuditHandler.Handle(c.Request.Context(), query)
	if err != nil {
		httperr.Respond(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AuthHandler) auditUserEvent(c *gin.Context, eventType, resourceID string, success bool, details map[string]string) {
	if h.auditLogger == nil {
		return
	}

	var actorID *uuid.UUID
	if rawActorID, exists := c.Get("userID"); exists {
		if parsed, err := uuid.Parse(rawActorID.(string)); err == nil {
			actorID = &parsed
		}
	}

	h.auditLogger.Log(c.Request.Context(), auditdomain.Event{
		EventType:    eventType,
		UserID:       actorID,
		ResourceType: "user",
		ResourceID:   resourceID,
		IPAddress:    c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
		Details:      details,
		Success:      success,
	})
}
