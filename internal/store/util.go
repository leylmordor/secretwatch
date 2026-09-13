package store

import (
	"path/filepath"
	"time"

	appconfig "github.com/leylmordor/secretwatch/internal/config"
)

func applyPolicyOverride(sec *Secret, overrides []appconfig.PolicyOverride) {
	for _, o := range overrides {
		if matched, _ := filepath.Match(o.Pattern, sec.Name); matched {
			sec.MaxAgeDays = o.MaxAgeDays
			return
		}
	}
}

func parseRFC3339(s string) (*time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
