package handlers

import (
	"net/http"
	"pr-reviewer/internal/service"

	"github.com/gin-gonic/gin"
)

type createPRRequest struct {
	ID       string `json:"pull_request_id" binding:"required"`
	Name     string `json:"pull_request_name" binding:"required"`
	AuthorID string `json:"author_id" binding:"required"`
}

func (h *Handler) createPR(c *gin.Context) {
	var req createPRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "INVALID_INPUT", "invalid input body")
		return
	}

	pr, err := h.svc.CreatePR(c.Request.Context(), req.ID, req.Name, req.AuthorID)
	if err != nil {
		switch err {
		case service.ErrPRExists:
			newErrorResponse(c, http.StatusConflict, "PR_EXISTS", "PR id already exists")
		case service.ErrAuthorNotFound:
			newErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "author or team not found")
		default:
			newErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		}
		return
	}

	reviewerIDs := make([]string, len(pr.Reviewers))
	for i, r := range pr.Reviewers {
		reviewerIDs[i] = r.ID
	}

	c.JSON(http.StatusCreated, gin.H{
		"pr": gin.H{
			"pull_request_id":    pr.ID,
			"pull_request_name":  pr.Name,
			"author_id":          pr.AuthorID,
			"status":             pr.Status,
			"assigned_reviewers": reviewerIDs,
		},
	})
}

type mergePRRequest struct {
	ID string `json:"pull_request_id" binding:"required"`
}

func (h *Handler) mergePR(c *gin.Context) {
	var req mergePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "INVALID_INPUT", "pull_request_id is required")
		return
	}

	pr, err := h.svc.MergePR(c.Request.Context(), req.ID)
	if err != nil {
		if err == service.ErrPRNotFound {
			newErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "pull request not found")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	reviewerIDs := make([]string, len(pr.Reviewers))
	for i, r := range pr.Reviewers {
		reviewerIDs[i] = r.ID
	}

	c.JSON(http.StatusOK, gin.H{
		"pr": gin.H{
			"pull_request_id":    pr.ID,
			"pull_request_name":  pr.Name,
			"author_id":          pr.AuthorID,
			"status":             pr.Status,
			"assigned_reviewers": reviewerIDs,
			"mergedAt":           pr.MergedAt,
		},
	})
}

type reassignReviewerRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
	OldUserID     string `json:"old_user_id" binding:"required"`
}

func (h *Handler) reassignReviewer(c *gin.Context) {
	var req reassignReviewerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "INVALID_INPUT", "invalid input body")
		return
	}

	pr, newReviewer, err := h.svc.ReassignReviewer(c.Request.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		switch err {
		case service.ErrPRNotFound:
			newErrorResponse(c, http.StatusNotFound, "NOT_FOUND", "pull request not found")
		case service.ErrPRMerged:
			newErrorResponse(c, http.StatusConflict, "PR_MERGED", "cannot reassign on merged PR")
		case service.ErrNotAssigned:
			newErrorResponse(c, http.StatusConflict, "NOT_ASSIGNED", "reviewer is not assigned to this PR")
		case service.ErrNoCandidate:
			newErrorResponse(c, http.StatusConflict, "NO_CANDIDATE", "no active replacement candidate in team")
		default:
			newErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		}
		return
	}

	reviewerIDs := make([]string, len(pr.Reviewers))
	for i, r := range pr.Reviewers {
		reviewerIDs[i] = r.ID
	}

	c.JSON(http.StatusOK, gin.H{
		"pr": gin.H{
			"pull_request_id":    pr.ID,
			"pull_request_name":  pr.Name,
			"author_id":          pr.AuthorID,
			"status":             pr.Status,
			"assigned_reviewers": reviewerIDs,
		},
		"replaced_by": newReviewer.ID,
	})
}

func (h *Handler) getReview(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		newErrorResponse(c, http.StatusBadRequest, "INVALID_INPUT", "user_id query param is required")
		return
	}

	prs, err := h.svc.GetUserReviews(c.Request.Context(), userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":       userID,
		"pull_requests": prs,
	})
}
