package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	appconfig "github.com/leylmordor/secretwatch/internal/config"
)

type VaultStore struct {
	cfg    appconfig.VaultConfig
	policy appconfig.RotationPolicyConfig
}

func NewVault(cfg appconfig.VaultConfig, policy appconfig.RotationPolicyConfig) *VaultStore {
	return &VaultStore{cfg: cfg, policy: policy}
}

func (s *VaultStore) Type() string { return "vault" }

func (s *VaultStore) List(ctx context.Context) ([]*Secret, error) {
	vcfg := api.DefaultConfig()
	vcfg.Address = s.cfg.Address

	client, err := api.NewClient(vcfg)
	if err != nil {
		return nil, fmt.Errorf("vault client: %w", err)
	}

	token := s.cfg.Token
	if token == "" {
		token = client.Token()
	}
	client.SetToken(token)

	if s.cfg.Namespace != "" {
		client.SetNamespace(s.cfg.Namespace)
	}

	var secrets []*Secret
	for _, path := range s.cfg.Paths {
		found, err := s.listPath(ctx, client, path)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, found...)
	}

	return secrets, nil
}

func (s *VaultStore) listPath(ctx context.Context, client *api.Client, path string) ([]*Secret, error) {
	// KV v2: mount/metadata/path
	metaPath := toMetadataPath(path)

	resp, err := client.Logical().ListWithContext(ctx, metaPath)
	if err != nil {
		return nil, fmt.Errorf("vault list %s: %w", path, err)
	}
	if resp == nil {
		return nil, nil
	}

	keys, ok := resp.Data["keys"].([]interface{})
	if !ok {
		return nil, nil
	}

	var secrets []*Secret
	for _, k := range keys {
		key := fmt.Sprintf("%v", k)
		fullPath := strings.TrimSuffix(path, "/") + "/" + key

		if strings.HasSuffix(key, "/") {
			// recurse into sub-path
			sub, err := s.listPath(ctx, client, fullPath)
			if err != nil {
				return nil, err
			}
			secrets = append(secrets, sub...)
			continue
		}

		sec := &Secret{
			Name:       key,
			FullPath:   fullPath,
			StoreType:  "vault",
			MaxAgeDays: s.policy.DefaultMaxAgeDays,
			Tags:       make(map[string]string),
		}

		meta, err := client.Logical().ReadWithContext(ctx, toMetadataPath(fullPath))
		if err == nil && meta != nil {
			if updated, ok := meta.Data["updated_time"].(string); ok {
				t, err := time.Parse(time.RFC3339Nano, updated)
				if err == nil {
					sec.LastRotated = &t
				}
			}
		}

		applyPolicyOverride(sec, s.policy.Overrides)
		sec.ComputeStatus()
		secrets = append(secrets, sec)
	}

	return secrets, nil
}

func toMetadataPath(path string) string {
	// Convert secret/foo -> secret/metadata/foo
	parts := strings.SplitN(strings.TrimPrefix(path, "/"), "/", 2)
	if len(parts) == 2 {
		return parts[0] + "/metadata/" + parts[1]
	}
	return path + "/metadata"
}
