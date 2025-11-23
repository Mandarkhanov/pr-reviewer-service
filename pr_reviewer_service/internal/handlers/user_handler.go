package handlers

import (
	"net/http"
	"pr-reviewer/internal/service"

	"github.com/gin-gonic/gin"
)

type setIsActiveRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	IsActive *bool  `json:"is_active" binding:"required"`
}

func (h *Handler) setIsActive(c *gin.Context) {
	var req setIsActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		return
	}

	user, err := h.svc.SetUserActive(c.Request.Context(), req.UserID, *req.IsActive)
	if err != nil {
		if err == service.ErrUserNotFound {
			newErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "user not found")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
