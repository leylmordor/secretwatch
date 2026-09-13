package history

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/logging/v2"
	"google.golang.org/api/option"
	appconfig "github.com/leylmordor/secretwatch/internal/config"
)

type GCPFetcher struct {
	cfg appconfig.GCPProjectConfig
}

func NewGCP(cfg appconfig.GCPProjectConfig) *GCPFetcher {
	return &GCPFetcher{cfg: cfg}
}

func (f *GCPFetcher) StoreType() string { return "gcp" }

func (f *GCPFetcher) Fetch(ctx context.Context, secretName string, limit int) ([]*Entry, error) {
	svc, err := logging.NewService(ctx, option.WithScopes(logging.LoggingReadScope))
	if err != nil {
		return nil, fmt.Errorf("gcp logging client: %w", err)
	}

	// Data Access audit logs must be enabled for Secret Manager in the GCP project.
	// Method names that represent a rotation/write:
	//   google.cloud.secretmanager.v1.SecretManagerService.AddSecretVersion
	//   google.cloud.secretmanager.v1.SecretManagerService.CreateSecret
	filter := fmt.Sprintf(
		`resource.type="secretmanager.googleapis.com/Secret" `+
			`AND resource.labels.secret_id="%s" `+
			`AND (protoPayload.methodName="google.cloud.secretmanager.v1.SecretManagerService.AddSecretVersion" `+
			`OR protoPayload.methodName="google.cloud.secretmanager.v1.SecretManagerService.CreateSecret")`,
		secretName,
	)

	req := &logging.ListLogEntriesRequest{
		ResourceNames: []string{fmt.Sprintf("projects/%s", f.cfg.Project)},
		Filter:        filter,
		OrderBy:       "timestamp desc",
		PageSize:      int64(limit),
	}

	resp, err := svc.Entries.List(req).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("gcp list log entries: %w", err)
	}

	var entries []*Entry
	for _, le := range resp.Entries {
		entry := &Entry{
			SecretName: secretName,
			Source:     "gcp-audit",
		}

		// ProtoPayload is raw JSON — unmarshal to extract method and actor.
		if le.ProtoPayload != nil {
			var payload map[string]interface{}
			if err := json.Unmarshal(le.ProtoPayload, &payload); err == nil {
				if method, ok := payload["methodName"].(string); ok {
					entry.Action = shortMethodName(method)
				}
				if auth, ok := payload["authenticationInfo"].(map[string]interface{}); ok {
					if principal, ok := auth["principalEmail"].(string); ok {
						entry.Actor = principal
					}
				}
			}
		}

		if le.Timestamp != "" {
			t, err := time.Parse(time.RFC3339Nano, le.Timestamp)
			if err == nil {
				entry.Timestamp = t
			}
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func shortMethodName(full string) string {
	// google.cloud.secretmanager.v1.SecretManagerService.AddSecretVersion → AddSecretVersion
	parts := splitDot(full)
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return full
}

func splitDot(s string) []string {
	var parts []string
	start := 0
	for i, c := range s {
		if c == '.' {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}
