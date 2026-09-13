package web

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	appconfig "github.com/leylmordor/secretwatch/internal/config"
	"github.com/leylmordor/secretwatch/internal/db"
	"github.com/leylmordor/secretwatch/internal/history"
	"github.com/leylmordor/secretwatch/internal/store"
)

type Server struct {
	db            *db.DB
	cfg           *appconfig.Config
	host          string
	port          int
	tmpl          *template.Template
	histTmpl      *template.Template
	secretsTmpl   *template.Template
	providersTmpl *template.Template
	reportsTmpl   *template.Template
	settingsTmpl  *template.Template
}

func New(database *db.DB, cfg *appconfig.Config, port int) (*Server, error) {
	funcs := template.FuncMap{
		"not":             func(b bool) bool { return !b },
		"statusClass":     statusClass,
		"formatTime":      formatTime,
		"formatTimeShort": formatTimeShort,
		"daysLabel":       daysLabel,
		"barPct":          barPct,
		"storeLabel":      storeLabel,
		"storeTabLabel":   storeTabLabel,
		"shortName":       shortName,
		"timeAgo":         timeAgo,
		"storeInitial":    storeInitial,
		"consoleURL":      consoleURL,
	}

	parse := func(name, src string) (*template.Template, error) {
		return template.New(name).Funcs(funcs).Parse(src)
	}

	tmpl, err := parse("dashboard", dashboardHTML)
	if err != nil {
		return nil, fmt.Errorf("dashboard template: %w", err)
	}
	histTmpl, err := parse("history", historyHTML)
	if err != nil {
		return nil, fmt.Errorf("history template: %w", err)
	}
	secretsTmpl, err := parse("secrets", secretsHTML)
	if err != nil {
		return nil, fmt.Errorf("secrets template: %w", err)
	}
	providersTmpl, err := parse("providers", providersHTML)
	if err != nil {
		return nil, fmt.Errorf("providers template: %w", err)
	}
	reportsTmpl, err := parse("reports", reportsHTML)
	if err != nil {
		return nil, fmt.Errorf("reports template: %w", err)
	}
	settingsTmpl, err := parse("settings", settingsHTML)
	if err != nil {
		return nil, fmt.Errorf("settings template: %w", err)
	}

	return &Server{
		db:            database,
		cfg:           cfg,
		host:          cfg.Web.Host,
		port:          port,
		tmpl:          tmpl,
		histTmpl:      histTmpl,
		secretsTmpl:   secretsTmpl,
		providersTmpl: providersTmpl,
		reportsTmpl:   reportsTmpl,
		settingsTmpl:  settingsTmpl,
	}, nil
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleDashboard)
	mux.HandleFunc("/history", s.handleHistory)
	mux.HandleFunc("/secrets", s.handleSecrets)
	mux.HandleFunc("/providers", s.handleProviders)
	mux.HandleFunc("/reports", s.handleReports)
	mux.HandleFunc("/settings", s.handleSettings)
	mux.HandleFunc("/api/secrets", s.handleAPISecrets)
	mux.HandleFunc("/api/history", s.handleAPIHistory)
	mux.HandleFunc("/api/exclude", s.handleAPIExclude)

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	fmt.Printf("Web UI available at http://%s\n", addr)
	return http.ListenAndServe(addr, mux)
}

// ── shared helpers ──────────────────────────────────────────────────────────

func (s *Server) baseCounts() (map[string]int, int) {
	all, _ := s.db.List("", "", "", false)
	counts := map[string]int{"ok": 0, "warning": 0, "overdue": 0, "unknown": 0, "total": len(all)}
	for _, sec := range all {
		counts[string(sec.Status)]++
	}
	allWithEx, _ := s.db.List("", "", "", true)
	return counts, len(allWithEx) - len(all)
}

func (s *Server) complianceRate(counts map[string]int) string {
	denom := counts["total"] - counts["unknown"]
	if denom < 1 {
		denom = 1
	}
	return fmt.Sprintf("%.1f%%", float64(counts["ok"])/float64(denom)*100)
}

func (s *Server) lastScanStr() string {
	if t, err := s.db.LastScan(); err == nil && t != nil {
		return t.Format("2006-01-02 15:04:05 UTC")
	}
	return "never"
}

// ── Dashboard ───────────────────────────────────────────────────────────────

type dashboardData struct {
	CurrentPage     string
	Secrets         []*store.Secret
	LastScan        string
	StoreFilter     string
	AccountFilter   string
	StatusFilter    string
	ShowExcluded    bool
	Counts          map[string]int
	ExcludedCount   int
	StoreTypes      []string
	StorePairs      []db.StoreAccountPair
	ComplianceRate  string
	RecentRotations []*store.Secret
	ActiveProviders int
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	storeFilter := r.URL.Query().Get("store")
	statusFilter := r.URL.Query().Get("status")
	showExcluded := r.URL.Query().Get("excluded") == "1"

	accountFilter := r.URL.Query().Get("account")
	secrets, err := s.db.List(storeFilter, accountFilter, statusFilter, showExcluded)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	counts, excludedCount := s.baseCounts()
	recentRotations, _ := s.db.RecentRotations(6)
	storeTypes := s.db.StoreTypes()
	storePairs, _ := s.db.StoreAccountPairs()

	data := dashboardData{
		CurrentPage:     "dashboard",
		Secrets:         secrets,
		LastScan:        s.lastScanStr(),
		StoreFilter:     storeFilter,
		AccountFilter:   accountFilter,
		StatusFilter:    statusFilter,
		ShowExcluded:    showExcluded,
		Counts:          counts,
		ExcludedCount:   excludedCount,
		StoreTypes:      storeTypes,
		StorePairs:      storePairs,
		ComplianceRate:  s.complianceRate(counts),
		RecentRotations: recentRotations,
		ActiveProviders: len(storeTypes),
	}

	w.Header().Set("Content-Type", "text/html")
	_ = s.tmpl.Execute(w, data)
}

// ── Secrets page ─────────────────────────────────────────────────────────

type secretsPageData struct {
	CurrentPage   string
	Secrets       []*store.Secret
	LastScan      string
	StoreFilter   string
	AccountFilter string
	StatusFilter  string
	ShowExcluded  bool
	Counts        map[string]int
	ExcludedCount int
	StoreTypes    []string
	StorePairs    []db.StoreAccountPair
}

func (s *Server) handleSecrets(w http.ResponseWriter, r *http.Request) {
	storeFilter := r.URL.Query().Get("store")
	accountFilter := r.URL.Query().Get("account")
	statusFilter := r.URL.Query().Get("status")
	showExcluded := r.URL.Query().Get("excluded") == "1"

	secrets, err := s.db.List(storeFilter, accountFilter, statusFilter, showExcluded)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	counts, excludedCount := s.baseCounts()
	storeTypes := s.db.StoreTypes()
	storePairs, _ := s.db.StoreAccountPairs()

	data := secretsPageData{
		CurrentPage:   "secrets",
		Secrets:       secrets,
		LastScan:      s.lastScanStr(),
		StoreFilter:   storeFilter,
		AccountFilter: accountFilter,
		StatusFilter:  statusFilter,
		ShowExcluded:  showExcluded,
		Counts:        counts,
		ExcludedCount: excludedCount,
		StoreTypes:    storeTypes,
		StorePairs:    storePairs,
	}

	w.Header().Set("Content-Type", "text/html")
	_ = s.secretsTmpl.Execute(w, data)
}

// ── Providers page ───────────────────────────────────────────────────────

type providersPageData struct {
	CurrentPage string
	Stats       []db.StoreStat
	LastScan    string
}

func (s *Server) handleProviders(w http.ResponseWriter, r *http.Request) {
	stats, _ := s.db.StoreStats()
	data := providersPageData{
		CurrentPage: "providers",
		Stats:       stats,
		LastScan:    s.lastScanStr(),
	}
	w.Header().Set("Content-Type", "text/html")
	_ = s.providersTmpl.Execute(w, data)
}

// ── Reports page ─────────────────────────────────────────────────────────

type reportsPageData struct {
	CurrentPage    string
	ComplianceRate string
	Counts         map[string]int
	Stats          []db.StoreStat
	LastScan       string
}

func (s *Server) handleReports(w http.ResponseWriter, r *http.Request) {
	counts, _ := s.baseCounts()
	stats, _ := s.db.StoreStats()
	data := reportsPageData{
		CurrentPage:    "reports",
		ComplianceRate: s.complianceRate(counts),
		Counts:         counts,
		Stats:          stats,
		LastScan:       s.lastScanStr(),
	}
	w.Header().Set("Content-Type", "text/html")
	_ = s.reportsTmpl.Execute(w, data)
}

// ── Settings page ────────────────────────────────────────────────────────

type settingsPageData struct {
	CurrentPage string
	Cfg         *appconfig.Config
}

func (s *Server) handleSettings(w http.ResponseWriter, r *http.Request) {
	data := settingsPageData{
		CurrentPage: "settings",
		Cfg:         s.cfg,
	}
	w.Header().Set("Content-Type", "text/html")
	_ = s.settingsTmpl.Execute(w, data)
}

// ── History ──────────────────────────────────────────────────────────────

type historyPageData struct {
	SecretName string
	StoreType  string
	Entries    []*history.Entry
	Error      string
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	secretName := r.URL.Query().Get("name")
	if secretName == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data := historyPageData{
		SecretName: secretName,
		StoreType:  s.db.GetStoreType(secretName),
	}
	entries, err := s.fetchHistory(r.Context(), secretName)
	if err != nil {
		data.Error = err.Error()
	} else {
		data.Entries = entries
	}

	w.Header().Set("Content-Type", "text/html")
	_ = s.histTmpl.Execute(w, data)
}

// ── API ───────────────────────────────────────────────────────────────────

func (s *Server) handleAPIHistory(w http.ResponseWriter, r *http.Request) {
	secretName := r.URL.Query().Get("name")
	if secretName == "" {
		http.Error(w, "name required", 400)
		return
	}
	entries, err := s.fetchHistory(r.Context(), secretName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(entries)
}

func (s *Server) fetchHistory(ctx context.Context, secretName string) ([]*history.Entry, error) {
	var fetchers []history.Fetcher
	for _, a := range s.cfg.Stores.AWS {
		fetchers = append(fetchers, history.NewAWS(a))
	}
	for _, p := range s.cfg.Stores.GCP {
		fetchers = append(fetchers, history.NewGCP(p))
	}
	if s.cfg.Stores.GitHub.Enabled {
		fetchers = append(fetchers, history.NewGitHub(s.cfg.Stores.GitHub))
	}

	var all []*history.Entry
	for _, f := range fetchers {
		entries, err := f.Fetch(ctx, secretName, 50)
		if err != nil {
			continue
		}
		all = append(all, entries...)
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].Timestamp.After(all[j].Timestamp)
	})
	return all, nil
}

func (s *Server) handleAPISecrets(w http.ResponseWriter, r *http.Request) {
	storeFilter := r.URL.Query().Get("store")
	statusFilter := r.URL.Query().Get("status")

	secrets, err := s.db.List(storeFilter, "", statusFilter, false)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(secrets)
}

func (s *Server) handleAPIExclude(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}
	fullPath := r.FormValue("full_path")
	excluded := r.FormValue("excluded") == "1"
	if fullPath == "" {
		http.Error(w, "full_path required", 400)
		return
	}
	if err := s.db.SetExcluded(fullPath, excluded); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── Template helpers ─────────────────────────────────────────────────────

var gcpHashRe = regexp.MustCompile(`-[0-9a-f]{8}-`)

func shortName(name string) string {
	loc := gcpHashRe.FindStringIndex(name)
	if loc == nil {
		return name
	}
	if suffix := name[loc[1]:]; suffix != "" {
		return suffix
	}
	return name
}

func storeLabel(t string) string {
	switch t {
	case "aws":
		return "AWS Secret"
	case "aws-iam":
		return "AWS IAM"
	case "gcp":
		return "GCP Secret"
	case "gcp-sa":
		return "Service Account"
	case "vault":
		return "Vault"
	case "github":
		return "GitHub"
	default:
		return t
	}
}

func storeTabLabel(t string) string {
	switch t {
	case "aws":
		return "AWS Secrets"
	case "aws-iam":
		return "AWS IAM"
	case "gcp":
		return "GCP Secrets"
	case "gcp-sa":
		return "Service Accounts"
	case "vault":
		return "Vault"
	case "github":
		return "GitHub"
	default:
		return t
	}
}

func barPct(daysUntil *int, maxAge int) int {
	if daysUntil == nil || maxAge == 0 {
		return 0
	}
	used := maxAge - *daysUntil
	if used < 0 {
		used = 0
	}
	pct := (used * 100) / maxAge
	if pct > 100 {
		pct = 100
	}
	return pct
}

func statusClass(s store.Status) string {
	switch s {
	case store.StatusOK:
		return "status-ok"
	case store.StatusWarning:
		return "status-warning"
	case store.StatusOverdue:
		return "status-overdue"
	default:
		return "status-unknown"
	}
}

func formatTime(t *time.Time) string {
	if t == nil {
		return "never"
	}
	return t.Format("Jan 02, 2006")
}

func formatTimeShort(t *time.Time) string {
	if t == nil {
		return "—"
	}
	return t.Format("Jan 02, 2006")
}

func timeAgo(t *time.Time) string {
	if t == nil {
		return "never"
	}
	d := time.Since(*t)
	if d < time.Hour {
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	}
	return fmt.Sprintf("%dd ago", int(d.Hours()/24))
}

func storeInitial(t string) string {
	switch t {
	case "aws", "aws-iam":
		return "A"
	case "gcp", "gcp-sa":
		return "G"
	case "vault":
		return "V"
	case "github":
		return "H"
	default:
		if len(t) > 0 {
			return string([]rune(t)[0])
		}
		return "?"
	}
}

func consoleURL(storeType, fullPath, region, account string) string {
	switch storeType {
	case "aws":
		// fullPath is the ARN: arn:aws:secretsmanager:us-east-1:123456:secret/name-XXXXX
		if region == "" {
			region = "us-east-1"
		}
		name := fullPath
		// extract secret name from ARN
		if idx := strings.LastIndex(fullPath, ":secret/"); idx >= 0 {
			name = fullPath[idx+len(":secret/"):]
		}
		return fmt.Sprintf("https://%s.console.aws.amazon.com/secretsmanager/secret?name=%s&region=%s", region, name, region)
	case "aws-iam":
		// fullPath: username/ACCESSKEYID
		username := fullPath
		if idx := strings.Index(fullPath, "/"); idx >= 0 {
			username = fullPath[:idx]
		}
		return fmt.Sprintf("https://us-east-1.console.aws.amazon.com/iamv2/home#/users/details/%s", username)
	case "gcp":
		// fullPath: projects/{project}/secrets/{name}
		project := account
		name := fullPath
		if idx := strings.LastIndex(fullPath, "/secrets/"); idx >= 0 {
			name = fullPath[idx+len("/secrets/"):]
		}
		return fmt.Sprintf("https://console.cloud.google.com/security/secret-manager/secret/%s/versions?project=%s", name, project)
	case "gcp-sa":
		// fullPath: sa-email/keyid
		saEmail := fullPath
		if idx := strings.Index(fullPath, "/"); idx >= 0 {
			saEmail = fullPath[:idx]
		}
		return fmt.Sprintf("https://console.cloud.google.com/iam-admin/serviceaccounts/details/%s/keys?project=%s", saEmail, account)
	case "github":
		if account == "" {
			return ""
		}
		return fmt.Sprintf("https://github.com/organizations/%s/settings/secrets/actions", account)
	case "vault":
		return "" // no standard console URL
	}
	return ""
}

func daysLabel(days *int) string {
	if days == nil {
		return "—"
	}
	if *days < 0 {
		return fmt.Sprintf("%d days overdue", -*days)
	}
	return fmt.Sprintf("%d days left", *days)
}
