package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Autodetect builds a Config from whatever credentials are available in the
// environment. Called when no config file exists.
func Autodetect() (*Config, error) {
	cfg := &Config{}
	cfg.applyDefaults()

	var detected []string

	if awsAccounts := detectAWS(); len(awsAccounts) > 0 {
		cfg.Stores.AWS = awsAccounts
		regions := strings.Join(awsAccounts[0].Regions, ", ")
		detected = append(detected, fmt.Sprintf("aws (%s)", regions))
	}

	if iamAccounts := detectAWSIAM(); len(iamAccounts) > 0 {
		cfg.Stores.AWSIAM = iamAccounts
		detected = append(detected, "aws_iam")
	}

	if gcpProjects := detectGCP(); len(gcpProjects) > 0 {
		cfg.Stores.GCP = gcpProjects
		cfg.Stores.GCPSAKeys = gcpProjects
		detected = append(detected, fmt.Sprintf("gcp (project: %s)", gcpProjects[0].Project))
	}

	if vaultCfg, ok := detectVault(); ok {
		cfg.Stores.Vault = vaultCfg
		detected = append(detected, fmt.Sprintf("vault (%s)", vaultCfg.Address))
	}

	if githubCfg, ok := detectGitHub(); ok {
		cfg.Stores.GitHub = githubCfg
		detected = append(detected, fmt.Sprintf("github (org: %s)", githubCfg.Org))
	}

	if len(detected) == 0 {
		return nil, fmt.Errorf("no credentials detected — run `secretwatch --help` or create ~/.secretwatch/config.yaml")
	}

	fmt.Fprintf(os.Stderr, "No config file found. Auto-detected: %s\n", strings.Join(detected, ", "))
	return cfg, nil
}

func detectAWS() []AWSAccountConfig {
	if !hasAWSCredentials() {
		return nil
	}

	region := firstNonEmpty(
		os.Getenv("AWS_DEFAULT_REGION"),
		os.Getenv("AWS_REGION"),
		"us-east-1",
	)

	return []AWSAccountConfig{{
		Profile: os.Getenv("AWS_PROFILE"),
		Regions: []string{region},
	}}
}

func detectAWSIAM() []AWSIAMAccountConfig {
	if !hasAWSCredentials() {
		return nil
	}
	return []AWSIAMAccountConfig{{
		Profile: os.Getenv("AWS_PROFILE"),
	}}
}

func hasAWSCredentials() bool {
	if os.Getenv("AWS_ACCESS_KEY_ID") != "" {
		return true
	}
	if os.Getenv("AWS_PROFILE") != "" {
		return true
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	for _, f := range []string{
		filepath.Join(home, ".aws", "credentials"),
		filepath.Join(home, ".aws", "config"),
	} {
		if _, err := os.Stat(f); err == nil {
			return true
		}
	}
	return false
}

func detectGCP() []GCPProjectConfig {
	if !hasGCPCredentials() {
		return nil
	}

	project := firstNonEmpty(
		os.Getenv("GOOGLE_CLOUD_PROJECT"),
		os.Getenv("GCLOUD_PROJECT"),
		os.Getenv("CLOUDSDK_CORE_PROJECT"),
		gcloudProject(),
	)
	if project == "" {
		fmt.Fprintf(os.Stderr, "GCP credentials detected but no project found. "+
			"Set GOOGLE_CLOUD_PROJECT or run: gcloud config set project <project-id>\n")
		return nil
	}

	return []GCPProjectConfig{{Project: project}}
}

func hasGCPCredentials() bool {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") != "" {
		return true
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	adc := filepath.Join(home, ".config", "gcloud", "application_default_credentials.json")
	_, err = os.Stat(adc)
	return err == nil
}

func gcloudProject() string {
	out, err := exec.Command("gcloud", "config", "get-value", "project", "--quiet").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func detectVault() (VaultConfig, bool) {
	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			data, err := os.ReadFile(filepath.Join(home, ".vault-token"))
			if err == nil {
				token = strings.TrimSpace(string(data))
			}
		}
	}
	if token == "" {
		return VaultConfig{}, false
	}

	addr := firstNonEmpty(os.Getenv("VAULT_ADDR"), "http://127.0.0.1:8200")

	return VaultConfig{
		Enabled: true,
		Address: addr,
		Token:   token,
		Paths:   []string{"secret/"},
	}, true
}

func detectGitHub() (GitHubConfig, bool) {
	token := firstNonEmpty(os.Getenv("GITHUB_TOKEN"), os.Getenv("GH_TOKEN"))
	org := os.Getenv("GITHUB_ORG")
	if token == "" || org == "" {
		return GitHubConfig{}, false
	}
	return GitHubConfig{Enabled: true, Org: org, Token: token}, true
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
