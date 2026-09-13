package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/leylmordor/secretwatch/internal/alert"
	appconfig "github.com/leylmordor/secretwatch/internal/config"
	"github.com/leylmordor/secretwatch/internal/db"
	"github.com/leylmordor/secretwatch/internal/history"
	"github.com/leylmordor/secretwatch/internal/scanner"
	"github.com/leylmordor/secretwatch/internal/store"
	"github.com/leylmordor/secretwatch/internal/web"
)

var configPath string

func main() {
	root := &cobra.Command{
		Use:   "secretwatch",
		Short: "Track secret rotation health across AWS, Vault, GCP, and GitHub",
	}
	root.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config file (default: ~/.secretwatch/config.yaml)")

	root.AddCommand(scanCmd(), listCmd(), serveCmd(), historyCmd(), alertCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func loadConfig() (*appconfig.Config, error) {
	return appconfig.Load(configPath)
}

func openDB(cfg *appconfig.Config) (*db.DB, error) {
	return db.Open(cfg.DBPath)
}

func buildStores(cfg *appconfig.Config) []store.Store {
	var stores []store.Store
	for _, a := range cfg.Stores.AWS {
		stores = append(stores, store.NewAWS(a, cfg.RotationPolicy))
	}
	for _, a := range cfg.Stores.AWSIAM {
		stores = append(stores, store.NewAWSIAM(a, cfg.RotationPolicy))
	}
	if cfg.Stores.Vault.Enabled {
		stores = append(stores, store.NewVault(cfg.Stores.Vault, cfg.RotationPolicy))
	}
	for _, p := range cfg.Stores.GCP {
		stores = append(stores, store.NewGCP(p, cfg.RotationPolicy))
	}
	for _, p := range cfg.Stores.GCPSAKeys {
		stores = append(stores, store.NewGCPServiceAccount(p, cfg.RotationPolicy))
	}
	if cfg.Stores.GitHub.Enabled {
		stores = append(stores, store.NewGitHub(cfg.Stores.GitHub, cfg.RotationPolicy))
	}
	return stores
}

func alertCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "alert",
		Short: "Send a Slack notification for overdue and expiring-soon secrets",
		Long: `Reads the cached scan results and sends a Slack alert for any secrets
that are overdue or expiring soon. Designed to run from cron after 'secretwatch scan'.

Example crontab (scan at 8am, alert at 8:05am daily):
  5 8 * * * secretwatch scan && secretwatch alert`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return fmt.Errorf("config: %w", err)
			}

			if cfg.Alerts.SlackWebhookURL == "" && !dryRun {
				return fmt.Errorf("alerts.slack_webhook_url not set in config")
			}

			notifyOn := cfg.Alerts.NotifyOn
			if len(notifyOn) == 0 {
				notifyOn = []string{"overdue", "warning"}
			}

			database, err := openDB(cfg)
			if err != nil {
				return fmt.Errorf("db: %w", err)
			}
			defer database.Close()

			var overdue, warning []*store.Secret
			for _, status := range notifyOn {
				secrets, err := database.List("", "", status, false)
				if err != nil {
					return fmt.Errorf("list %s: %w", status, err)
				}
				switch status {
				case "overdue":
					overdue = secrets
				case "warning":
					warning = secrets
				}
			}

			fmt.Printf("Found %d overdue, %d expiring soon.\n", len(overdue), len(warning))

			if len(overdue) == 0 && len(warning) == 0 {
				fmt.Println("Nothing to alert on.")
				return nil
			}

			if dryRun {
				fmt.Println("Dry run — would send Slack alert:")
				for _, s := range overdue {
					fmt.Printf("  [overdue] %s (%s)\n", s.Name, s.StoreType)
				}
				for _, s := range warning {
					fmt.Printf("  [warning] %s (%s)\n", s.Name, s.StoreType)
				}
				return nil
			}

			notifier := alert.NewSlack(cfg.Alerts.SlackWebhookURL)
			if err := notifier.Send(overdue, warning); err != nil {
				return fmt.Errorf("slack: %w", err)
			}

			fmt.Println("Slack alert sent.")
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print what would be alerted without sending")
	return cmd
}

func scanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "scan",
		Short: "Scan all configured secret stores and cache results",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return fmt.Errorf("config: %w", err)
			}

			database, err := openDB(cfg)
			if err != nil {
				return fmt.Errorf("db: %w", err)
			}
			defer database.Close()

			stores := buildStores(cfg)
			if len(stores) == 0 {
				return fmt.Errorf("no stores enabled in config")
			}

			fmt.Printf("Scanning %d store(s)...\n", len(stores))
			result := scanner.Run(context.Background(), stores)

			for _, err := range result.Errors {
				fmt.Fprintf(os.Stderr, "warning: %v\n", err)
			}

			if err := database.Save(result.Secrets); err != nil {
				return fmt.Errorf("save: %w", err)
			}

			fmt.Printf("Found %d secrets.\n", len(result.Secrets))
			printTable(result.Secrets)
			return nil
		},
	}
}

func listCmd() *cobra.Command {
	var storeFilter, statusFilter string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List cached secrets from last scan",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return fmt.Errorf("config: %w", err)
			}

			database, err := openDB(cfg)
			if err != nil {
				return fmt.Errorf("db: %w", err)
			}
			defer database.Close()

			secrets, err := database.List(storeFilter, "", statusFilter, false)
			if err != nil {
				return fmt.Errorf("list: %w", err)
			}

			if len(secrets) == 0 {
				fmt.Println("No secrets found. Run `secretwatch scan` first.")
				return nil
			}

			printTable(secrets)
			return nil
		},
	}

	cmd.Flags().StringVar(&storeFilter, "store", "", "filter by store type (aws, vault, gcp, github)")
	cmd.Flags().StringVar(&statusFilter, "status", "", "filter by status (ok, warning, overdue, unknown)")
	return cmd
}

func serveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the web dashboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig()
			if err != nil {
				return fmt.Errorf("config: %w", err)
			}

			database, err := openDB(cfg)
			if err != nil {
				return fmt.Errorf("db: %w", err)
			}
			defer database.Close()

			srv, err := web.New(database, cfg, cfg.Web.Port)
			if err != nil {
				return err
			}

			return srv.Start()
		},
	}
}

func historyCmd() *cobra.Command {
	var storeType string
	var limit int

	cmd := &cobra.Command{
		Use:   "history <secret-name>",
		Short: "Show who rotated a secret and when (from audit logs)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			secretName := args[0]

			cfg, err := loadConfig()
			if err != nil {
				return fmt.Errorf("config: %w", err)
			}

			fetchers := buildFetchers(cfg, storeType)
			if len(fetchers) == 0 {
				return fmt.Errorf("no matching stores enabled in config")
			}

			var all []*history.Entry
			for _, f := range fetchers {
				entries, err := f.Fetch(context.Background(), secretName, limit)
				if err != nil {
					fmt.Fprintf(os.Stderr, "warning [%s]: %v\n", f.StoreType(), err)
					continue
				}
				all = append(all, entries...)
			}

			if len(all) == 0 {
				fmt.Println("No rotation history found.")
				fmt.Println("Note: AWS requires CloudTrail; GCP requires Data Access audit logs to be enabled.")
				return nil
			}

			sort.Slice(all, func(i, j int) bool {
				return all[i].Timestamp.After(all[j].Timestamp)
			})

			printHistory(all)
			return nil
		},
	}

	cmd.Flags().StringVar(&storeType, "store", "", "limit to a specific store (aws, gcp, github)")
	cmd.Flags().IntVar(&limit, "limit", 20, "max number of events to fetch per store")
	return cmd
}

func buildFetchers(cfg *appconfig.Config, storeType string) []history.Fetcher {
	var fetchers []history.Fetcher
	if storeType == "" || storeType == "aws" {
		for _, a := range cfg.Stores.AWS {
			fetchers = append(fetchers, history.NewAWS(a))
		}
	}
	if storeType == "" || storeType == "gcp" {
		for _, p := range cfg.Stores.GCP {
			fetchers = append(fetchers, history.NewGCP(p))
		}
	}
	if (storeType == "" || storeType == "github") && cfg.Stores.GitHub.Enabled {
		fetchers = append(fetchers, history.NewGitHub(cfg.Stores.GitHub))
	}
	return fetchers
}

func printHistory(entries []*history.Entry) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tACTOR\tACTION\tSTORE\tNOTE")
	fmt.Fprintln(w, "---------\t-----\t------\t-----\t----")
	for _, e := range entries {
		ts := e.Timestamp.Format("2006-01-02 15:04:05")
		region := e.Region
		storeLabel := e.Source
		if region != "" {
			storeLabel += "/" + region
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", ts, e.Actor, e.Action, storeLabel, e.Note)
	}
	w.Flush()
}

func printTable(secrets []*store.Secret) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "STATUS\tNAME\tSTORE\tLAST ROTATED\tEXPIRY")
	fmt.Fprintln(w, "------\t----\t-----\t------------\t------")

	for _, s := range secrets {
		lastRotated := "never"
		if s.LastRotated != nil {
			lastRotated = s.LastRotated.Format("2006-01-02")
		}
		expiry := "—"
		if s.DaysUntilExpiry != nil {
			if *s.DaysUntilExpiry < 0 {
				expiry = fmt.Sprintf("%d days overdue", -*s.DaysUntilExpiry)
			} else {
				expiry = fmt.Sprintf("%d days left", *s.DaysUntilExpiry)
			}
		}
		storeLabel := s.StoreType
		if s.Region != "" {
			storeLabel += "/" + s.Region
		}

		status := statusIcon(s.Status) + " " + string(s.Status)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", status, s.Name, storeLabel, lastRotated, expiry)
	}

	w.Flush()
}

func statusIcon(s store.Status) string {
	switch s {
	case store.StatusOK:
		return "✓"
	case store.StatusWarning:
		return "!"
	case store.StatusOverdue:
		return "✗"
	default:
		return "?"
	}
}
