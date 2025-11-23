package handlers

import (
	"net/http"
	"pr-reviewer/internal/service"

	"github.com/gin-gonic/gin"
)

type createTeamRequest struct {
	TeamName string `json:"team_name" binding:"required"`
	Members  []struct {
		UserID   string `json:"user_id" binding:"required"`
		Username string `json:"username" binding:"required"`
		IsActive bool   `json:"is_active" binding:"required"`
	} `json:"members"`
}

func (h *Handler) createTeam(c *gin.Context) {
	var req createTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "INVALID_INPUT", "invalid input body")
		return
	}
	team := toDomainTeam(req)

	err := h.svc.CreateTeam(c.Request.Context(), team)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "TEAM_EXISTS", err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"team": team})
}

func (h *Handler) getTeam(c *gin.Context) {
	name := c.Query("team_name")
	if name == "" {
		newErrorResponse(c, http.StatusBadRequest, "INVALID_INPUT", "team_name required")
		return
	}

	team, err := h.svc.GetTeam(c.Request.Context(), name)
	if err != nil {
		if err == service.ErrTeamNotFound {
			newErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "team not found")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, team)
}
