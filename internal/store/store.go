package store

import (
	"context"
	"time"
)

type Status string

const (
	StatusOK      Status = "ok"
	StatusWarning Status = "warning"
	StatusOverdue Status = "overdue"
	StatusUnknown Status = "unknown"
)

type Secret struct {
	Name            string
	FullPath        string
	StoreType       string
	Region          string
	Account         string
	LastRotated     *time.Time
	MaxAgeDays      int
	Status          Status
	DaysUntilExpiry *int
	Tags            map[string]string
	Excluded        bool
}

func (s *Secret) ComputeStatus() {
	if s.LastRotated == nil {
		s.Status = StatusUnknown
		return
	}

	age := int(time.Since(*s.LastRotated).Hours() / 24)
	remaining := s.MaxAgeDays - age
	s.DaysUntilExpiry = &remaining

	switch {
	case remaining < 0:
		s.Status = StatusOverdue
	case remaining <= 14:
		s.Status = StatusWarning
	default:
		s.Status = StatusOK
	}
}

type Store interface {
	Type() string
	List(ctx context.Context) ([]*Secret, error)
}
