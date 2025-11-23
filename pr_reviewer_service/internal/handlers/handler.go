package handlers

import (
	"pr-reviewer/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) InitRoutes(router *gin.Engine) {
	router.POST("/team/add", h.createTeam)
	router.GET("/team/get", h.getTeam)

	router.POST("/users/setIsActive", h.setIsActive)
	router.GET("/users/getReview", h.getReview)

	router.POST("/pullRequest/create", h.createPR)
	router.POST("/pullRequest/merge", h.mergePR)
	router.POST("/pullRequest/reassign", h.reassignReviewer)
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, code string, msg string) {
	c.AbortWithStatusJSON(statusCode, errorResponse{
		Error: errorDetail{
			Code:    code,
			Message: msg,
		},
	})
}
