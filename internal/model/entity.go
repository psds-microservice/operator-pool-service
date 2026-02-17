package model

import (
	"time"

	"github.com/google/uuid"
)

type OperatorStatus struct {
	UserID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	Available      bool      `gorm:"not null;default:false" json:"available"`
	ActiveSessions int       `gorm:"not null;default:0" json:"active_sessions"`
	MaxSessions    int       `gorm:"not null;default:5" json:"max_sessions"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (OperatorStatus) TableName() string { return "operator_status" }
