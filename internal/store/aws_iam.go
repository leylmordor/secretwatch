package store

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	appconfig "github.com/leylmordor/secretwatch/internal/config"
)

// AWSIAMStore scans IAM access keys for all users in the account.
// These keys have no expiry date but should be rotated on a schedule.
type AWSIAMStore struct {
	cfg    appconfig.AWSIAMAccountConfig
	policy appconfig.RotationPolicyConfig
}

func NewAWSIAM(cfg appconfig.AWSIAMAccountConfig, policy appconfig.RotationPolicyConfig) *AWSIAMStore {
	return &AWSIAMStore{cfg: cfg, policy: policy}
}

func (s *AWSIAMStore) Type() string { return "aws-iam" }

func (s *AWSIAMStore) List(ctx context.Context) ([]*Secret, error) {
	opts := []func(*config.LoadOptions) error{}
	if s.cfg.Profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(s.cfg.Profile))
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("aws iam config: %w", err)
	}

	client := iam.NewFromConfig(awsCfg)
	var secrets []*Secret

	paginator := iam.NewListUsersPaginator(client, &iam.ListUsersInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("iam list users: %w", err)
		}

		for _, user := range page.Users {
			username := aws.ToString(user.UserName)

			keyResp, err := client.ListAccessKeys(ctx, &iam.ListAccessKeysInput{
				UserName: user.UserName,
			})
			if err != nil {
				continue
			}

			account := s.cfg.Profile
			if account == "" {
				account = "default"
			}
			for _, key := range keyResp.AccessKeyMetadata {
				if key.Status == "Inactive" {
					continue
				}

				name := fmt.Sprintf("%s/%s", username, aws.ToString(key.AccessKeyId))
				sec := &Secret{
					Name:       name,
					FullPath:   fmt.Sprintf("iam/users/%s/access-keys/%s", username, aws.ToString(key.AccessKeyId)),
					StoreType:  "aws-iam",
					Account:    account,
					MaxAgeDays: s.policy.DefaultMaxAgeDays,
					Tags:       map[string]string{"user": username, "key_id": aws.ToString(key.AccessKeyId)},
				}

				if key.CreateDate != nil {
					sec.LastRotated = key.CreateDate
				}

				applyPolicyOverride(sec, s.policy.Overrides)
				sec.ComputeStatus()
				secrets = append(secrets, sec)
			}
		}
	}

	return secrets, nil
}
