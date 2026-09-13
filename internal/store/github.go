package store

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
	appconfig "github.com/leylmordor/secretwatch/internal/config"
)

type GitHubStore struct {
	cfg    appconfig.GitHubConfig
	policy appconfig.RotationPolicyConfig
}

func NewGitHub(cfg appconfig.GitHubConfig, policy appconfig.RotationPolicyConfig) *GitHubStore {
	return &GitHubStore{cfg: cfg, policy: policy}
}

func (s *GitHubStore) Type() string { return "github" }

func (s *GitHubStore) List(ctx context.Context) ([]*Secret, error) {
	token := s.cfg.Token
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("github: no token provided (set token in config or GITHUB_TOKEN env)")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	var secrets []*Secret
	opt := &github.ListOptions{PerPage: 100}

	for {
		page, resp, err := client.Actions.ListOrgSecrets(ctx, s.cfg.Org, opt)
		if err != nil {
			return nil, fmt.Errorf("github list org secrets: %w", err)
		}

		for _, gs := range page.Secrets {
			sec := &Secret{
				Name:       gs.Name,
				FullPath:   fmt.Sprintf("orgs/%s/actions/secrets/%s", s.cfg.Org, gs.Name),
				StoreType:  "github",
				Account:    s.cfg.Org,
				MaxAgeDays: s.policy.DefaultMaxAgeDays,
				Tags:       make(map[string]string),
			}

			// GitHub exposes UpdatedAt on secrets
			if !gs.UpdatedAt.IsZero() {
				t := gs.UpdatedAt.Time
				sec.LastRotated = &t
			}

			applyPolicyOverride(sec, s.policy.Overrides)
			sec.ComputeStatus()
			secrets = append(secrets, sec)
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return secrets, nil
}
