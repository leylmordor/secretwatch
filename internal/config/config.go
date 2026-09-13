package config

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Stores         StoresConfig         `yaml:"stores"`
	RotationPolicy RotationPolicyConfig `yaml:"rotation_policy"`
	DBPath         string               `yaml:"db_path"`
	Web            WebConfig            `yaml:"web"`
	Alerts         AlertsConfig         `yaml:"alerts"`
}

type StoresConfig struct {
	AWS       []AWSAccountConfig    `yaml:"aws"`
	AWSIAM    []AWSIAMAccountConfig `yaml:"aws_iam"`
	Vault     VaultConfig           `yaml:"vault"`
	GCP       []GCPProjectConfig    `yaml:"gcp"`
	GCPSAKeys []GCPProjectConfig    `yaml:"gcp_sa_keys"`
	GitHub    GitHubConfig          `yaml:"github"`
}

type AWSAccountConfig struct {
	Profile string   `yaml:"profile"`
	Regions []string `yaml:"regions"`
	Label   string   `yaml:"label"`
}

type AWSIAMAccountConfig struct {
	Profile string `yaml:"profile"`
	Label   string `yaml:"label"`
}

type GCPProjectConfig struct {
	Project string `yaml:"project"`
	Label   string `yaml:"label"`
}

type VaultConfig struct {
	Enabled   bool     `yaml:"enabled"`
	Address   string   `yaml:"address"`
	Token     string   `yaml:"token"`
	Paths     []string `yaml:"paths"`
	Namespace string   `yaml:"namespace"`
}

type GitHubConfig struct {
	Enabled bool   `yaml:"enabled"`
	Org     string `yaml:"org"`
	Token   string `yaml:"token"`
}

type RotationPolicyConfig struct {
	DefaultMaxAgeDays int              `yaml:"default_max_age_days"`
	Overrides         []PolicyOverride `yaml:"overrides"`
}

type PolicyOverride struct {
	Pattern    string `yaml:"pattern"`
	MaxAgeDays int    `yaml:"max_age_days"`
}

type WebConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type AlertsConfig struct {
	SlackWebhookURL string   `yaml:"slack_webhook_url"`
	NotifyOn        []string `yaml:"notify_on"`
}

func Load(path string) (*Config, error) {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(home, ".secretwatch", "config.yaml")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Autodetect()
		}
		return nil, err
	}

	data = []byte(os.ExpandEnv(string(data)))

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	cfg.applyDefaults()
	return cfg, nil
}

func (c *Config) applyDefaults() {
	if c.RotationPolicy.DefaultMaxAgeDays == 0 {
		c.RotationPolicy.DefaultMaxAgeDays = 90
	}
	if c.DBPath == "" {
		home, _ := os.UserHomeDir()
		c.DBPath = filepath.Join(home, ".secretwatch", "cache.db")
	} else if strings.HasPrefix(c.DBPath, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			c.DBPath = filepath.Join(home, c.DBPath[2:])
		}
	}
	if c.Web.Port == 0 {
		c.Web.Port = 8080
	}
	if c.Web.Host == "" {
		if h := os.Getenv("SECRETWATCH_WEB_HOST"); h != "" {
			c.Web.Host = h
		} else {
			c.Web.Host = "127.0.0.1"
		}
	}
	if c.Stores.Vault.Address == "" {
		c.Stores.Vault.Address = "http://127.0.0.1:8200"
	}
	for i := range c.Stores.AWS {
		if len(c.Stores.AWS[i].Regions) == 0 {
			c.Stores.AWS[i].Regions = []string{"us-east-1"}
		}
	}
}
