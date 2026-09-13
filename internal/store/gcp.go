package store

import (
	"context"
	"fmt"
	"strings"

	secretmanager "google.golang.org/api/secretmanager/v1"
	"google.golang.org/api/option"
	appconfig "github.com/leylmordor/secretwatch/internal/config"
)

type GCPStore struct {
	cfg    appconfig.GCPProjectConfig
	policy appconfig.RotationPolicyConfig
}

func NewGCP(cfg appconfig.GCPProjectConfig, policy appconfig.RotationPolicyConfig) *GCPStore {
	return &GCPStore{cfg: cfg, policy: policy}
}

func (s *GCPStore) Type() string { return "gcp" }

func (s *GCPStore) List(ctx context.Context) ([]*Secret, error) {
	svc, err := secretmanager.NewService(ctx, option.WithScopes(secretmanager.CloudPlatformScope))
	if err != nil {
		return nil, fmt.Errorf("gcp secret manager client: %w", err)
	}

	var secrets []*Secret
	parent := fmt.Sprintf("projects/%s", s.cfg.Project)

	err = svc.Projects.Secrets.List(parent).Context(ctx).Pages(ctx, func(page *secretmanager.ListSecretsResponse) error {
		for _, s2 := range page.Secrets {
			name := secretName(s2.Name)
			sec := &Secret{
				Name:       name,
				FullPath:   s2.Name,
				StoreType:  "gcp",
				Account:    s.cfg.Project,
				MaxAgeDays: s.policy.DefaultMaxAgeDays,
				Tags:       make(map[string]string),
			}

			for k, v := range s2.Labels {
				sec.Tags[k] = v
			}

			// GCP doesn't expose last rotation directly on the secret resource;
			// we use the latest version's create time as a proxy.
			latest, err := svc.Projects.Secrets.Versions.Get(s2.Name + "/versions/latest").Context(ctx).Do()
			if err == nil && latest.CreateTime != "" {
				// CreateTime is RFC3339
				var t interface{}
				_ = t
				if parsed, err := parseRFC3339(latest.CreateTime); err == nil {
					sec.LastRotated = parsed
				}
			}

			applyPolicyOverride(sec, s.policy.Overrides)
			sec.ComputeStatus()
			secrets = append(secrets, sec)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("gcp list secrets: %w", err)
	}

	return secrets, nil
}

func secretName(fullName string) string {
	parts := strings.Split(fullName, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return fullName
}
