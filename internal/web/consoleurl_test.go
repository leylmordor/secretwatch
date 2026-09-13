package web

import (
	"strings"
	"testing"
)

func TestConsoleURL_AWS(t *testing.T) {
	url := consoleURL("aws", "arn:aws:secretsmanager:us-east-1:123:secret/prod-db", "us-east-1", "123")
	if !strings.Contains(url, "us-east-1.console.aws.amazon.com") {
		t.Errorf("expected AWS console URL, got %q", url)
	}
	if !strings.Contains(url, "prod-db") {
		t.Errorf("expected secret name in URL, got %q", url)
	}
}

func TestConsoleURL_AWSIAM(t *testing.T) {
	url := consoleURL("aws-iam", "alice/AKIAIOSFODNN7EXAMPLE", "", "")
	if !strings.Contains(url, "alice") {
		t.Errorf("expected username in URL, got %q", url)
	}
	if !strings.Contains(url, "iamv2") {
		t.Errorf("expected IAM URL, got %q", url)
	}
}

func TestConsoleURL_GCP(t *testing.T) {
	url := consoleURL("gcp", "projects/myproject/secrets/my-secret", "", "myproject")
	if !strings.Contains(url, "my-secret") {
		t.Errorf("expected secret name in URL, got %q", url)
	}
	if !strings.Contains(url, "myproject") {
		t.Errorf("expected project in URL, got %q", url)
	}
}

func TestConsoleURL_GitHub(t *testing.T) {
	url := consoleURL("github", "orgs/myorg/actions/secrets/TOKEN", "", "myorg")
	if !strings.Contains(url, "myorg") {
		t.Errorf("expected org in URL, got %q", url)
	}
}

func TestConsoleURL_GitHubEmptyAccount(t *testing.T) {
	// Verify the bug is fixed: empty account produces a useless URL
	url := consoleURL("github", "orgs//actions/secrets/TOKEN", "", "")
	if strings.Contains(url, "organizations//") {
		t.Errorf("GitHub URL has empty org segment: %q", url)
	}
}

func TestConsoleURL_Vault(t *testing.T) {
	url := consoleURL("vault", "secret/prod/db", "", "")
	if url != "" {
		t.Errorf("expected empty URL for vault, got %q", url)
	}
}
