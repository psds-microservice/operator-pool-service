package service

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/psds-microservice/operator-pool-service/internal/errs"
	"github.com/psds-microservice/operator-pool-service/internal/model"
	"gorm.io/gorm"
)

type OperatorService struct {
	db  *gorm.DB
	mu  sync.Mutex
	idx int // round-robin index
}

func NewOperatorService(db *gorm.DB) *OperatorService {
	return &OperatorService{db: db}
}

func (s *OperatorService) SetStatus(userID uuid.UUID, available bool, maxSessions int) error {
	if maxSessions <= 0 {
		maxSessions = 5
	}
	var op model.OperatorStatus
	s.db.Where("user_id = ?", userID).FirstOrCreate(&op, model.OperatorStatus{UserID: userID, MaxSessions: maxSessions})
	return s.db.Model(&op).Updates(map[string]interface{}{
		"available":    available,
		"max_sessions": maxSessions,
		"updated_at":   time.Now(),
	}).Error
}

func (s *OperatorService) UpdateStatus(userID uuid.UUID, available bool) error {
	return s.db.Model(&model.OperatorStatus{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
		"available":  available,
		"updated_at": time.Now(),
	}).Error
}

func (s *OperatorService) IncrementSessions(userID uuid.UUID) error {
	return s.db.Model(&model.OperatorStatus{}).Where("user_id = ?", userID).UpdateColumn("active_sessions", gorm.Expr("active_sessions + ?", 1)).Error
}

func (s *OperatorService) DecrementSessions(userID uuid.UUID) error {
	return s.db.Model(&model.OperatorStatus{}).Where("user_id = ? AND active_sessions > 0", userID).UpdateColumn("active_sessions", gorm.Expr("active_sessions - ?", 1)).Error
}

// Next returns the next available operator (round-robin).
func (s *OperatorService) Next() (uuid.UUID, error) {
	var list []model.OperatorStatus
	if err := s.db.Where("available = ? AND active_sessions < max_sessions", true).Find(&list).Error; err != nil {
		return uuid.Nil, err
	}
	if len(list) == 0 {
		return uuid.Nil, errs.ErrNoOperatorAvailable
	}
	s.mu.Lock()
	idx := s.idx % len(list)
	s.idx++
	s.mu.Unlock()
	return list[idx].UserID, nil
}

func (s *OperatorService) Stats() (available int, total int, err error) {
	var avail int64
	s.db.Model(&model.OperatorStatus{}).Where("available = ?", true).Count(&avail)
	var tot int64
	s.db.Model(&model.OperatorStatus{}).Count(&tot)
	return int(avail), int(tot), nil
}

func (s *OperatorService) ListAvailable() ([]model.OperatorStatus, error) {
	var list []model.OperatorStatus
	err := s.db.Where("available = ? AND active_sessions < max_sessions", true).Find(&list).Error
	return list, err
}

// ListAll returns all operators with their status (for directory/search).
func (s *OperatorService) ListAll() ([]model.OperatorStatus, error) {
	var list []model.OperatorStatus
	err := s.db.Find(&list).Error
	return list, err
}
