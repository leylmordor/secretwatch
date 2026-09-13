package store

import (
	"testing"
	"time"
)

func TestComputeStatus_NilLastRotated(t *testing.T) {
	s := &Secret{MaxAgeDays: 90}
	s.ComputeStatus()
	if s.Status != StatusUnknown {
		t.Errorf("expected unknown, got %s", s.Status)
	}
	if s.DaysUntilExpiry != nil {
		t.Errorf("expected nil DaysUntilExpiry, got %d", *s.DaysUntilExpiry)
	}
}

func TestComputeStatus_OK(t *testing.T) {
	t30 := time.Now().AddDate(0, 0, -30)
	s := &Secret{MaxAgeDays: 90, LastRotated: &t30}
	s.ComputeStatus()
	if s.Status != StatusOK {
		t.Errorf("expected ok, got %s", s.Status)
	}
	if s.DaysUntilExpiry == nil || *s.DaysUntilExpiry != 60 {
		v := -1
		if s.DaysUntilExpiry != nil {
			v = *s.DaysUntilExpiry
		}
		t.Errorf("expected DaysUntilExpiry 60, got %d", v)
	}
}

func TestComputeStatus_Warning(t *testing.T) {
	t80 := time.Now().AddDate(0, 0, -80)
	s := &Secret{MaxAgeDays: 90, LastRotated: &t80}
	s.ComputeStatus()
	if s.Status != StatusWarning {
		t.Errorf("expected warning, got %s", s.Status)
	}
}

func TestComputeStatus_Overdue(t *testing.T) {
	t100 := time.Now().AddDate(0, 0, -100)
	s := &Secret{MaxAgeDays: 90, LastRotated: &t100}
	s.ComputeStatus()
	if s.Status != StatusOverdue {
		t.Errorf("expected overdue, got %s", s.Status)
	}
	if s.DaysUntilExpiry == nil || *s.DaysUntilExpiry >= 0 {
		t.Errorf("expected negative DaysUntilExpiry, got %v", s.DaysUntilExpiry)
	}
}

func TestComputeStatus_ExactlyAtBoundary(t *testing.T) {
	// exactly 14 days remaining => warning
	t76 := time.Now().AddDate(0, 0, -76)
	s := &Secret{MaxAgeDays: 90, LastRotated: &t76}
	s.ComputeStatus()
	if s.Status != StatusWarning {
		t.Errorf("expected warning at 14-day boundary, got %s", s.Status)
	}
}
