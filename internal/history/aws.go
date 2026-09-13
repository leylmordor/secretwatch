package history

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	appconfig "github.com/leylmordor/secretwatch/internal/config"
)

// rotationEvents are the CloudTrail event names that represent a secret being written/rotated.
var rotationEvents = map[string]bool{
	"RotateSecret":   true,
	"PutSecretValue": true,
	"UpdateSecret":   true,
	"CreateSecret":   true,
}

type AWSFetcher struct {
	cfg appconfig.AWSAccountConfig
}

func NewAWS(cfg appconfig.AWSAccountConfig) *AWSFetcher {
	return &AWSFetcher{cfg: cfg}
}

func (f *AWSFetcher) StoreType() string { return "aws" }

func (f *AWSFetcher) Fetch(ctx context.Context, secretName string, limit int) ([]*Entry, error) {
	var all []*Entry

	for _, region := range f.cfg.Regions {
		opts := []func(*config.LoadOptions) error{
			config.WithRegion(region),
		}
		if f.cfg.Profile != "" {
			opts = append(opts, config.WithSharedConfigProfile(f.cfg.Profile))
		}

		awsCfg, err := config.LoadDefaultConfig(ctx, opts...)
		if err != nil {
			return nil, fmt.Errorf("aws config (%s): %w", region, err)
		}

		client := cloudtrail.NewFromConfig(awsCfg)

		// CloudTrail lets us look up events by resource name (the secret name or ARN).
		paginator := cloudtrail.NewLookupEventsPaginator(client, &cloudtrail.LookupEventsInput{
			LookupAttributes: []types.LookupAttribute{
				{
					AttributeKey:   types.LookupAttributeKeyResourceName,
					AttributeValue: aws.String(secretName),
				},
			},
			MaxResults: aws.Int32(int32(limit)),
		})

		for paginator.HasMorePages() && len(all) < limit {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, fmt.Errorf("cloudtrail lookup (%s): %w", region, err)
			}

			for _, event := range page.Events {
				if !rotationEvents[aws.ToString(event.EventName)] {
					continue
				}

				entry := &Entry{
					SecretName: secretName,
					Action:     aws.ToString(event.EventName),
					Source:     "cloudtrail",
					Region:     region,
				}

				if event.EventTime != nil {
					entry.Timestamp = *event.EventTime
				}

				// Actor: prefer username, fall back to the event's user identity ARN.
				if event.Username != nil && *event.Username != "" {
					entry.Actor = *event.Username
				}

				// Distinguish automated rotation from manual updates.
				if aws.ToString(event.EventName) == "RotateSecret" {
					entry.Note = "automated rotation"
				} else if isLambdaActor(entry.Actor) {
					entry.Note = "rotation lambda"
				}

				all = append(all, entry)
			}
		}
	}

	return all, nil
}

func isLambdaActor(actor string) bool {
	return strings.Contains(actor, "lambda") || strings.Contains(actor, "Lambda")
}
