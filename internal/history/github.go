package history

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
	appconfig "github.com/leylmordor/secretwatch/internal/config"
)

// GitHub audit log actions related to secret writes.
var githubSecretActions = map[string]bool{
	"org.create_actions_secret": true,
	"org.update_actions_secret": true,
	"org.remove_actions_secret": true,
}

type GitHubFetcher struct {
	cfg appconfig.GitHubConfig
}

func NewGitHub(cfg appconfig.GitHubConfig) *GitHubFetcher {
	return &GitHubFetcher{cfg: cfg}
}

func (f *GitHubFetcher) StoreType() string { return "github" }

func (f *GitHubFetcher) Fetch(ctx context.Context, secretName string, limit int) ([]*Entry, error) {
	token := f.cfg.Token
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("github: no token provided")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	client := github.NewClient(oauth2.NewClient(ctx, ts))

	// GitHub audit log phrase search finds entries mentioning the secret name.
	phrase := secretName
	include := "all"
	opts := &github.GetAuditLogOptions{
		Phrase:  &phrase,
		Include: &include,
		ListCursorOptions: github.ListCursorOptions{
			PerPage: limit,
		},
	}

	auditEntries, _, err := client.Organizations.GetAuditLog(ctx, f.cfg.Org, opts)
	if err != nil {
		return nil, fmt.Errorf("github audit log: %w", err)
	}

	var entries []*Entry
	for _, ae := range auditEntries {
		action := ae.GetAction()
		if !githubSecretActions[action] {
			continue
		}
		// Extra filter: ensure this event is actually about our secret.
		// The secret name lives in AdditionalFields["name"] for org secret events.
		if name, ok := ae.AdditionalFields["name"].(string); !ok || !strings.EqualFold(name, secretName) {
			continue
		}

		entry := &Entry{
			SecretName: secretName,
			Action:     action,
			Actor:      ae.GetActor(),
			Source:     "github-audit",
		}

		if ts := ae.GetCreatedAt(); !ts.IsZero() {
			entry.Timestamp = ts.Time
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
