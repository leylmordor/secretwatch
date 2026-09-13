package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/leylmordor/secretwatch/internal/store"
)

type SlackNotifier struct {
	webhookURL string
}

func NewSlack(webhookURL string) *SlackNotifier {
	return &SlackNotifier{webhookURL: webhookURL}
}

type slackPayload struct {
	Text   string       `json:"text"`
	Blocks []slackBlock `json:"blocks"`
}

type slackBlock struct {
	Type string     `json:"type"`
	Text *slackText `json:"text,omitempty"`
}

type slackText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (n *SlackNotifier) Send(overdue, warning []*store.Secret) error {
	if len(overdue) == 0 && len(warning) == 0 {
		return nil
	}

	blocks := []slackBlock{
		{
			Type: "header",
			Text: &slackText{Type: "plain_text", Text: "🔐 SecretWatch — Rotation Alert"},
		},
	}

	if len(overdue) > 0 {
		lines := []string{fmt.Sprintf("*🔴 %d secret(s) overdue for rotation:*", len(overdue))}
		for _, s := range overdue {
			days := ""
			if s.DaysUntilExpiry != nil {
				days = fmt.Sprintf(" — %d days overdue", -*s.DaysUntilExpiry)
			}
			lines = append(lines, fmt.Sprintf("• `%s` (%s)%s", s.Name, storeLabel(s), days))
		}
		blocks = append(blocks, slackBlock{
			Type: "section",
			Text: &slackText{Type: "mrkdwn", Text: strings.Join(lines, "\n")},
		})
	}

	if len(warning) > 0 {
		lines := []string{fmt.Sprintf("*⚠️  %d secret(s) expiring soon:*", len(warning))}
		for _, s := range warning {
			days := ""
			if s.DaysUntilExpiry != nil {
				days = fmt.Sprintf(" — %d days left", *s.DaysUntilExpiry)
			}
			lines = append(lines, fmt.Sprintf("• `%s` (%s)%s", s.Name, storeLabel(s), days))
		}
		blocks = append(blocks, slackBlock{
			Type: "section",
			Text: &slackText{Type: "mrkdwn", Text: strings.Join(lines, "\n")},
		})
	}

	blocks = append(blocks, slackBlock{
		Type: "context",
		Text: &slackText{Type: "mrkdwn", Text: "Run `secretwatch list --status overdue` or open the dashboard to see details."},
	})

	payload := slackPayload{
		Text:   fmt.Sprintf("SecretWatch: %d overdue, %d expiring soon", len(overdue), len(warning)),
		Blocks: blocks,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack returned %d", resp.StatusCode)
	}

	return nil
}

func storeLabel(s *store.Secret) string {
	if s.Region != "" {
		return s.StoreType + "/" + s.Region
	}
	return s.StoreType
}
