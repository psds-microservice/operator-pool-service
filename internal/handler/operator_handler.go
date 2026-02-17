package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/psds-microservice/operator-pool-service/internal/errs"
	"github.com/psds-microservice/operator-pool-service/internal/service"
)

type OperatorHandler struct {
	svc *service.OperatorService
}

func NewOperatorHandler(svc *service.OperatorService) *OperatorHandler {
	return &OperatorHandler{svc: svc}
}

func (h *OperatorHandler) Status(c *gin.Context) {
	var req struct {
		UserID      string `json:"user_id"`
		Available   bool   `json:"available"`
		MaxSessions int    `json:"max_sessions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}
	if err := h.svc.SetStatus(userID, req.Available, req.MaxSessions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *OperatorHandler) Next(c *gin.Context) {
	operatorID, err := h.svc.Next()
	if err != nil {
		if errors.Is(err, errs.ErrNoOperatorAvailable) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no operator available"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"operator_id": operatorID.String()})
}

func (h *OperatorHandler) Stats(c *gin.Context) {
	available, total, err := h.svc.Stats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"available": available, "total": total})
}

// List returns all operators with their status (for operator-directory and other consumers).
func (h *OperatorHandler) List(c *gin.Context) {
	list, err := h.svc.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	out := make([]gin.H, 0, len(list))
	for _, op := range list {
		out = append(out, gin.H{
			"user_id":         op.UserID.String(),
			"available":       op.Available,
			"active_sessions": op.ActiveSessions,
			"max_sessions":    op.MaxSessions,
			"updated_at":      op.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"operators": out})
}
