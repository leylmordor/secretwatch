package store

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	appconfig "github.com/leylmordor/secretwatch/internal/config"
)

type AWSStore struct {
	cfg    appconfig.AWSAccountConfig
	policy appconfig.RotationPolicyConfig
}

func NewAWS(cfg appconfig.AWSAccountConfig, policy appconfig.RotationPolicyConfig) *AWSStore {
	return &AWSStore{cfg: cfg, policy: policy}
}

func (s *AWSStore) Type() string { return "aws" }

func (s *AWSStore) List(ctx context.Context) ([]*Secret, error) {
	var secrets []*Secret

	for _, region := range s.cfg.Regions {
		opts := []func(*config.LoadOptions) error{
			config.WithRegion(region),
		}
		if s.cfg.Profile != "" {
			opts = append(opts, config.WithSharedConfigProfile(s.cfg.Profile))
		}

		awsCfg, err := config.LoadDefaultConfig(ctx, opts...)
		if err != nil {
			return nil, fmt.Errorf("aws config (%s): %w", region, err)
		}

		client := secretsmanager.NewFromConfig(awsCfg)
		paginator := secretsmanager.NewListSecretsPaginator(client, &secretsmanager.ListSecretsInput{})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, fmt.Errorf("aws list secrets (%s): %w", region, err)
			}

			account := s.cfg.Profile
			if account == "" {
				account = "default"
			}
			for _, entry := range page.SecretList {
				sec := &Secret{
					Name:       aws.ToString(entry.Name),
					FullPath:   aws.ToString(entry.ARN),
					StoreType:  "aws",
					Region:     region,
					Account:    account,
					MaxAgeDays: s.policy.DefaultMaxAgeDays,
					Tags:       make(map[string]string),
				}

				if entry.LastRotatedDate != nil {
					sec.LastRotated = entry.LastRotatedDate
				} else if entry.LastChangedDate != nil {
					sec.LastRotated = entry.LastChangedDate
				}

				for _, tag := range entry.Tags {
					if tag.Key != nil && tag.Value != nil {
						sec.Tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
					}
				}

				applyPolicyOverride(sec, s.policy.Overrides)
				sec.ComputeStatus()
				secrets = append(secrets, sec)
			}
		}
	}

	return secrets, nil
}
