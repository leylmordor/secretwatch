package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	appconfig "github.com/leylmordor/secretwatch/internal/config"
)

// GCPServiceAccountStore scans GCP service account keys.
// SA keys have no enforced expiry but should be rotated on a schedule.
type GCPServiceAccountStore struct {
	cfg    appconfig.GCPProjectConfig
	policy appconfig.RotationPolicyConfig
}

func NewGCPServiceAccount(cfg appconfig.GCPProjectConfig, policy appconfig.RotationPolicyConfig) *GCPServiceAccountStore {
	return &GCPServiceAccountStore{cfg: cfg, policy: policy}
}

func (s *GCPServiceAccountStore) Type() string { return "gcp-sa" }

func (s *GCPServiceAccountStore) List(ctx context.Context) ([]*Secret, error) {
	svc, err := iam.NewService(ctx, option.WithScopes(iam.CloudPlatformScope))
	if err != nil {
		return nil, fmt.Errorf("gcp iam client: %w", err)
	}

	var secrets []*Secret
	project := fmt.Sprintf("projects/%s", s.cfg.Project)

	// List all service accounts in the project.
	err = svc.Projects.ServiceAccounts.List(project).Context(ctx).Pages(ctx,
		func(page *iam.ListServiceAccountsResponse) error {
			for _, sa := range page.Accounts {
				keys, err := svc.Projects.ServiceAccounts.Keys.List(sa.Name).
					KeyTypes("USER_MANAGED"). // skip Google-managed keys
					Context(ctx).Do()
				if err != nil {
					return nil // skip this SA if we can't list keys
				}

				for _, key := range keys.Keys {
					saEmail := sa.Email
					keyID := lastSegment(key.Name)

					sec := &Secret{
						Name:       fmt.Sprintf("%s/%s", saEmail, keyID[:min(len(keyID), 12)]),
						FullPath:   key.Name,
						StoreType:  "gcp-sa",
						Account:    s.cfg.Project,
						MaxAgeDays: s.policy.DefaultMaxAgeDays,
						Tags:       map[string]string{"service_account": saEmail, "key_id": keyID},
					}

					if key.ValidAfterTime != "" {
						t, err := time.Parse(time.RFC3339, key.ValidAfterTime)
						if err == nil {
							sec.LastRotated = &t
						}
					}

					applyPolicyOverride(sec, s.policy.Overrides)
					sec.ComputeStatus()
					secrets = append(secrets, sec)
				}
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("gcp list service accounts: %w", err)
	}

	return secrets, nil
}

func lastSegment(name string) string {
	parts := strings.Split(name, "/")
	return parts[len(parts)-1]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
