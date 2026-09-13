package web

const historyHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>SecretWatch — History</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Geist:wght@400;500;600&family=Geist+Mono:wght@400&display=swap" rel="stylesheet">
<style>
  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
  :root {
    --bg: #09090b; --surface: #0c0c0e; --border: #1f1f23;
    --text: #ededf0; --muted: #71717a; --accent: #4285f4;
    --overdue: #ef4444; --warning: #f59e0b; --ok: #10b981;
    --font: 'Geist', system-ui, sans-serif;
    --mono: 'Geist Mono', monospace;
  }
  body { background: var(--bg); color: var(--text); font-family: var(--font); font-size: 14px; line-height: 1.5; -webkit-font-smoothing: antialiased; }
  header { display: flex; align-items: center; gap: 20px; padding: 0 40px; height: 52px; border-bottom: 1px solid var(--border); }
  .logo { font-size: 14px; font-weight: 600; color: var(--text); }
  .logo span { color: var(--muted); font-weight: 400; }
  .back { font-size: 13px; color: var(--muted); text-decoration: none; }
  .back:hover { color: var(--text); }
  .content { padding: 40px; max-width: 1000px; }
  .secret-title { font-family: var(--mono); font-size: 18px; font-weight: 400; color: var(--text); margin-bottom: 4px; word-break: break-all; }
  .secret-short { font-size: 24px; font-weight: 600; margin-bottom: 6px; }
  .subtitle { font-size: 13px; color: var(--muted); margin-bottom: 32px; }
  .notice { border-left: 3px solid var(--border); padding: 10px 16px; font-size: 13px; color: var(--muted); margin-bottom: 28px; line-height: 1.6; }
  table { width: 100%; border-collapse: collapse; }
  thead th { text-align: left; padding: 8px 16px; font-size: 12px; font-weight: 500; color: var(--muted); border-bottom: 1px solid var(--border); }
  tbody tr { border-bottom: 1px solid var(--border); }
  tbody tr:hover td { background: rgba(255,255,255,0.02); }
  td { padding: 14px 16px; font-size: 13px; }
  .actor { font-family: var(--mono); font-size: 11px; color: var(--muted); max-width: 300px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .action { font-weight: 500; }
  .source-tag { font-size: 11px; background: rgba(255,255,255,0.06); border-radius: 4px; padding: 2px 8px; color: var(--muted); }
  .note { font-size: 12px; color: var(--muted); }
  .empty { text-align: center; padding: 80px 32px; color: var(--muted); }
  .empty p { margin-top: 8px; font-size: 13px; }
  .error { border-left: 3px solid var(--overdue); padding: 12px 16px; color: var(--overdue); margin-bottom: 24px; font-size: 13px; }
</style>
</head>
<body>
<header>
  <a href="/" class="back">← back</a>
  <div class="logo">SecretWatch <span>rotation monitor</span></div>
</header>
<div class="content">
  {{$short := shortName .SecretName}}
  {{if ne $short .SecretName}}
  <div class="secret-short">{{$short}}</div>
  <div class="secret-title">{{.SecretName}}</div>
  {{else}}
  <div class="secret-short">{{.SecretName}}</div>
  {{end}}
  <div class="subtitle">Rotation history from audit logs</div>

  {{if .Error}}
  <div class="error">{{.Error}}</div>
  {{else if .Entries}}
  <table>
    <thead>
      <tr>
        <th>When</th>
        <th>Actor</th>
        <th>Action</th>
        <th>Source</th>
        <th>Note</th>
      </tr>
    </thead>
    <tbody>
    {{range .Entries}}
    <tr>
      <td>{{.Timestamp.Format "2006-01-02 15:04:05"}}</td>
      <td class="actor">{{.Actor}}</td>
      <td class="action">{{.Action}}</td>
      <td><span class="source-tag">{{.Source}}{{if .Region}}/{{.Region}}{{end}}</span></td>
      <td class="note">{{.Note}}</td>
    </tr>
    {{end}}
    </tbody>
  </table>
  {{else}}
  {{if eq .StoreType "aws"}}
  <div class="notice">History is pulled from CloudTrail. Make sure CloudTrail is enabled in your account (it is by default). 90-day window.</div>
  {{else if eq .StoreType "gcp"}}
  <div class="notice">History requires Cloud Audit Logs with Data Access enabled for Secret Manager. That setting is off by default — enable it in IAM &amp; Admin › Audit Logs in your GCP project.</div>
  {{else if eq .StoreType "github"}}
  <div class="notice">History comes from the GitHub Audit Log API and requires an org admin token. Free plans retain 90 days.</div>
  {{else if eq .StoreType "vault"}}
  <div class="notice">Vault doesn't expose audit history via API. Logs are written to file or syslog on the Vault server.</div>
  {{else}}
  <div class="notice">No rotation history found. Make sure audit logging is enabled for this store.</div>
  {{end}}
  <div class="empty">
    <p>No rotation events found for this secret.</p>
  </div>
  {{end}}
</div>
</body>
</html>
`


const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>SecretWatch</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Geist:wght@400;500;600;700&family=Geist+Mono:wght@400;500;600;700&display=swap" rel="stylesheet">
<style>
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
:root {
  --bg:         #09090b;
  --sidebar-bg: #09090b;
  --card:       #0c0c0e;
  --card2:      #141416;
  --border:     #1f1f23;
  --text:       #ededf0;
  --sub:        #a1a1aa;
  --muted:      #71717a;
  --faint:      #52525b;
  --green:      #10b981;
  --green-bg:   rgba(16,185,129,0.1);
  --orange:     #f59e0b;
  --orange-bg:  rgba(245,158,11,0.1);
  --red:        #ef4444;
  --red-bg:     rgba(239,68,68,0.1);
  --blue:       #4285f4;
  --aws:        #ff9900;
  --font:       'Geist', system-ui, sans-serif;
  --mono:       'Geist Mono', monospace;
}
html, body { height: 100%; }
body { background: var(--bg); color: var(--text); font-family: var(--font); font-size: 13px; line-height: 1.5; -webkit-font-smoothing: antialiased; display: flex; height: 100vh; overflow: hidden; }

/* ── Sidebar ── */
.sidebar {
  width: 64px; flex-shrink: 0;
  background: var(--sidebar-bg);
  border-right: 1px solid var(--border);
  display: flex; flex-direction: column; align-items: center; justify-content: space-between;
  padding: 24px 0;
}
.sidebar-top { display: flex; flex-direction: column; align-items: center; gap: 32px; }
.sidebar-logo {
  width: 36px; height: 36px; background: var(--card);
  border: 1px solid var(--border);
  border-radius: 8px; display: flex; align-items: center; justify-content: center;
  font-size: 18px; font-weight: 700; color: var(--text); font-family: var(--mono);
}
.sidebar-nav { display: flex; flex-direction: column; gap: 12px; }
.sidebar-icon {
  width: 48px; height: 48px; border-radius: 8px;
  border: 1px solid transparent;
  display: flex; align-items: center; justify-content: center;
  color: var(--muted); cursor: pointer; transition: all 0.12s;
}
.sidebar-icon:hover, .sidebar-icon.active { background: var(--card); border-color: var(--border); color: var(--text); }
.sidebar-icon svg { width: 19px; height: 19px; }

/* ── Main ── */
.main { flex: 1; overflow-y: auto; display: flex; flex-direction: column; min-width: 0; }

/* Page header */
.page-header {
  display: flex; align-items: flex-start; justify-content: space-between;
  padding: 24px 28px 20px; flex-shrink: 0;
}
.page-header h1 { font-size: 20px; font-weight: 700; letter-spacing: -0.3px; }
.page-header .subtitle { font-size: 13px; color: var(--sub); margin-top: 2px; }
.header-search {
  display: flex; align-items: center; gap: 10px;
}
.search-box {
  position: relative;
}
.search-box svg { position: absolute; left: 10px; top: 50%; transform: translateY(-50%); color: var(--muted); pointer-events: none; }
#search {
  background: var(--card); border: 1px solid var(--border); border-radius: 8px;
  color: var(--text); font-family: var(--font); font-size: 13px;
  padding: 8px 14px 8px 32px; width: 240px; outline: none; transition: border-color 0.15s;
}
#search:focus { border-color: var(--blue); }
#search::placeholder { color: var(--muted); }

/* Stats row */
.stats-row { display: grid; grid-template-columns: repeat(4,1fr); gap: 12px; padding: 0 28px 20px; }
.stat-card {
  background: var(--card); border: 1px solid var(--border); border-radius: 10px;
  padding: 18px 20px; text-decoration: none; transition: border-color 0.12s;
}
.stat-card:hover { border-color: #2d3654; }
.stat-card-top { display: flex; align-items: center; justify-content: space-between; margin-bottom: 10px; }
.stat-card-label { font-size: 12px; color: var(--sub); font-weight: 500; }
.stat-tag {
  font-size: 10px; font-weight: 600; padding: 2px 6px; border-radius: 4px;
  font-family: var(--mono);
}
.stat-tag.orange { background: var(--orange-bg); color: var(--orange); }
.stat-tag.red    { background: var(--red-bg);    color: var(--red); }
.stat-tag.green  { background: var(--green-bg);  color: var(--green); }
.stat-n { font-size: 22px; font-weight: 700; font-family: var(--mono); font-variant-numeric: tabular-nums; color: var(--text); }
.stat-n.compliance { color: var(--green); }

/* Provider cards */
.section-header { padding: 4px 28px 10px; font-size: 12px; font-weight: 600; color: var(--muted); }
.provider-row { display: grid; grid-template-columns: repeat(4,1fr); gap: 8px; padding: 0 28px 20px; }
.provider-card {
  background: var(--card); border-radius: 6px;
  padding: 10px 12px; display: flex; align-items: center; gap: 10px;
}
.provider-icon {
  width: 24px; height: 24px; border-radius: 4px;
  display: flex; align-items: center; justify-content: center;
  font-size: 11px; font-weight: 700; flex-shrink: 0; font-family: var(--mono);
}
.provider-icon.aws    { background: rgba(255,153,0,0.15);   color: var(--aws); }
.provider-icon.gcp    { background: rgba(66,133,244,0.15);  color: var(--blue); }
.provider-icon.vault  { background: rgba(239,68,68,0.15);   color: var(--red); }
.provider-icon.github { background: rgba(161,161,170,0.15); color: var(--sub); }
.provider-icon.other  { background: rgba(113,113,122,0.15); color: var(--muted); }
.provider-name { font-size: 12px; font-weight: 600; color: var(--text); }
.provider-sync { font-size: 10px; color: var(--muted); margin-top: 1px; display: flex; align-items: center; gap: 4px; font-family: var(--mono); }
.sync-dot { width: 6px; height: 6px; border-radius: 50%; background: var(--green); }
.sync-dot.syncing { background: var(--orange); }

/* Table */
.table-section { padding: 0 28px 28px; }
.table-top { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px; }
.table-top h2 { font-size: 14px; font-weight: 600; }
.table-top-right { display: flex; gap: 8px; align-items: center; }
.filter-tabs { display: flex; gap: 2px; }
.ftab {
  font-size: 12px; font-weight: 500; color: var(--muted);
  padding: 4px 8px; border-radius: 4px; background: rgba(24,24,27,1); text-decoration: none; transition: all 0.12s;
}
.ftab:hover { color: var(--sub); }
.ftab.active { color: var(--sub); }
.excl-btn { font-size: 11px; color: var(--muted); background: none; border: none; cursor: pointer; font-family: var(--font); padding: 4px 8px; transition: color 0.1s; }
.excl-btn:hover { color: var(--sub); }
.provider-select {
  background: var(--card2); border: 1px solid var(--border); border-radius: 6px;
  color: var(--sub); font-family: var(--font); font-size: 12px;
  padding: 5px 28px 5px 10px; outline: none; cursor: pointer;
  appearance: none; -webkit-appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='10' height='6' viewBox='0 0 10 6'%3E%3Cpath d='M1 1l4 4 4-4' stroke='%2371717a' stroke-width='1.5' fill='none' stroke-linecap='round'/%3E%3C/svg%3E");
  background-repeat: no-repeat; background-position: right 8px center;
  min-width: 160px; transition: border-color 0.12s;
}
.provider-select:hover { border-color: var(--muted); color: var(--text); }
.provider-select option { background: #18181b; color: var(--text); }

table { width: 100%; border-collapse: collapse; background: var(--card); border-radius: 8px; overflow: hidden; }
thead th {
  text-align: left; font-size: 11px; font-weight: 600; color: var(--muted);
  padding: 10px 14px; border-bottom: 1px solid var(--border);
}
tbody tr { border-bottom: 1px solid var(--border); transition: background 0.1s; }
tbody tr:last-child { border-bottom: none; }
tbody tr:hover { background: rgba(24,24,27,0.6); }
tbody tr.excluded td { opacity: 0.35; }
td { padding: 10px 14px; vertical-align: middle; }

.td-name { font-family: var(--mono); font-size: 13px; color: var(--text); font-weight: 400; display: block; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 240px; }
.td-fullpath { font-family: var(--mono); font-size: 10px; color: var(--faint); display: block; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 240px; margin-top: 1px; }
.td-provider { font-size: 12px; color: var(--sub); display: flex; align-items: center; gap: 6px; }
.td-date { font-size: 12px; font-family: var(--mono); color: var(--sub); white-space: nowrap; }
.td-policy { font-size: 12px; font-family: var(--mono); color: var(--muted); white-space: nowrap; }
.td-status { white-space: nowrap; }

.status-badge {
  display: inline-flex; align-items: center;
  font-size: 11px; font-weight: 600; padding: 2px 8px; border-radius: 12px;
}
.status-badge.ok      { background: var(--green-bg);  color: var(--green); }
.status-badge.warning { background: var(--orange-bg); color: var(--orange); }
.status-badge.overdue { background: var(--red-bg);    color: var(--red); }
.status-badge.unknown { color: var(--muted); background: rgba(113,113,122,0.1); }
.status-badge.excluded { color: var(--muted); background: rgba(113,113,122,0.1); font-style: italic; }

.td-actions { text-align: right; width: 44px; }
.row-menu { position: relative; display: inline-block; }
.menu-btn {
  background: none; border: none; cursor: pointer; color: var(--muted);
  font-size: 18px; letter-spacing: 1px; padding: 2px 6px; border-radius: 5px;
  font-family: var(--font); transition: all 0.1s; line-height: 1;
}
.menu-btn:hover { color: var(--text); background: var(--border); }
.dropdown {
  display: none; position: absolute; right: 0; top: calc(100% + 4px);
  background: #1e2438; border: 1px solid #2a3250; border-radius: 8px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.5); min-width: 150px; z-index: 100; overflow: hidden;
}
.dropdown.open { display: block; }
.dd-item {
  display: block; width: 100%; text-align: left; padding: 9px 14px;
  font-size: 13px; color: var(--sub); background: none; border: none;
  cursor: pointer; font-family: var(--font); text-decoration: none; transition: all 0.1s;
}
.dd-item:hover { color: var(--text); background: rgba(255,255,255,0.05); }
.dd-item.is-excluded { color: var(--blue); }

.tbl-footer { padding: 10px 14px; border-top: 1px solid var(--border); font-size: 12px; color: var(--muted); background: var(--card2); display: flex; justify-content: space-between; align-items: center; }
.excl-link { color: var(--muted); text-decoration: none; }
.excl-link:hover { color: var(--sub); }

.empty-row td { text-align: center; padding: 48px 14px; color: var(--sub); }

/* ── Right panel ── */
.right-panel {
  width: 280px; flex-shrink: 0;
  border-left: 1px solid var(--border);
  overflow-y: auto;
  padding: 24px 0 0;
}
.rp-section { padding: 0 20px 24px; }
.rp-title { font-size: 14px; font-weight: 700; margin-bottom: 4px; }
.rp-subtitle { font-size: 11px; color: var(--muted); margin-bottom: 16px; }
.rp-divider { height: 1px; background: var(--border); margin: 0 20px 24px; }

.timeline-item { display: flex; align-items: flex-start; gap: 10px; margin-bottom: 16px; }
.timeline-item:last-child { margin-bottom: 0; }
.tl-icon {
  width: 28px; height: 28px; border-radius: 50%; flex-shrink: 0;
  display: flex; align-items: center; justify-content: center;
  font-size: 11px; font-weight: 700;
}
.tl-icon.aws    { background: rgba(245,147,50,0.15);  color: var(--orange); }
.tl-icon.gcp    { background: rgba(76,142,245,0.15);  color: var(--blue); }
.tl-icon.vault  { background: rgba(224,69,69,0.15);   color: var(--red); }
.tl-icon.github { background: rgba(122,133,160,0.15); color: var(--sub); }
.tl-icon.other  { background: var(--card2); color: var(--muted); }
.tl-body { flex: 1; min-width: 0; }
.tl-name { font-family: var(--mono); font-size: 11px; font-weight: 500; color: var(--text); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.tl-action { font-size: 11px; color: var(--muted); margin-top: 1px; }
.tl-time { font-size: 10px; color: var(--muted); flex-shrink: 0; }

.compliance-rate { font-size: 36px; font-weight: 700; color: var(--green); letter-spacing: -1px; margin-bottom: 4px; }
.compliance-label { font-size: 12px; color: var(--sub); margin-bottom: 10px; }
.compliance-note { font-size: 12px; color: var(--muted); line-height: 1.5; }
</style>
</head>
<body>

<!-- Sidebar -->
<aside class="sidebar">
  <div class="sidebar-top">
    <a href="/" class="sidebar-logo" style="text-decoration:none">S</a>
    <nav class="sidebar-nav">
      <a href="/" class="sidebar-icon {{if eq .CurrentPage "dashboard"}}active{{end}}" title="Overview">
        <svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M2.375 7.125H16.625M2.375 11.875H16.625M7.125 2.375V16.625M11.875 2.375V16.625M3.95833 2.375H15.0417C15.9161 2.375 16.625 3.08388 16.625 3.95833V15.0417C16.625 15.9161 15.9161 16.625 15.0417 16.625H3.95833C3.08388 16.625 2.375 15.9161 2.375 15.0417V3.95833C2.375 3.08388 3.08388 2.375 3.95833 2.375Z"/></svg>
      </a>
      <a href="/secrets" class="sidebar-icon {{if eq .CurrentPage "secrets"}}active{{end}}" title="Secrets">
        <svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M11.0831 5.54166L13.4582 7.91685M7.44128 8.39158L2.04664 13.7864C1.74967 14.0833 1.58279 14.486 1.5827 14.9059V16.6256C1.5827 16.8356 1.66611 17.0369 1.81458 17.1854C1.96306 17.3339 2.16443 17.4173 2.3744 17.4173H4.7495C4.95947 17.4173 5.16085 17.3339 5.30932 17.1854C5.45779 17.0369 5.5412 16.8356 5.5412 16.6256V15.8338C5.5412 15.6239 5.62461 15.4225 5.77309 15.274C5.92156 15.1255 6.12293 15.0421 6.3329 15.0421H7.1246C7.33458 15.0421 7.53595 14.9587 7.68442 14.8102C7.83289 14.6617 7.91631 14.4604 7.91631 14.2504V13.4587C7.91631 13.2487 7.99972 13.0473 8.14819 12.8988C8.29666 12.7503 8.49803 12.6669 8.70801 12.6669H8.84418C9.26409 12.6668 9.66677 12.4999 9.96364 12.203L10.6081 11.5585M9.81652 2.13736C10.1825 1.79792 10.6632 1.6093 11.1624 1.6093C11.6616 1.6093 12.1423 1.79792 12.5083 2.13736L16.8627 6.49187C17.2021 6.85788 17.3907 7.33863 17.3907 7.83781C17.3907 8.337 17.2021 8.81775 16.8627 9.18376L13.9334 12.1132C13.5674 12.4526 13.0866 12.6412 12.5875 12.6412C12.0883 12.6412 11.6076 12.4526 11.2416 12.1132L6.88723 7.75864C6.5478 7.39263 6.35919 6.91189 6.35919 6.4127C6.35919 5.91352 6.5478 5.43277 6.88723 5.06676L9.81652 2.13736Z"/></svg>
      </a>
      <a href="/providers" class="sidebar-icon {{if eq .CurrentPage "providers"}}active{{end}}" title="Providers">
        <svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M16.625 3.95789C16.625 5.26967 13.435 6.33308 9.5 6.33308C5.56497 6.33308 2.375 5.26967 2.375 3.95789M16.625 3.95789C16.625 2.64611 13.435 1.5827 9.5 1.5827C5.56497 1.5827 2.375 2.64611 2.375 3.95789M16.625 3.95789V15.0421C16.625 15.6721 15.8743 16.2762 14.5381 16.7216C13.2019 17.1671 11.3897 17.4173 9.5 17.4173C7.61033 17.4173 5.79806 17.1671 4.46186 16.7216C3.12567 16.2762 2.375 15.6721 2.375 15.0421V3.95789M2.375 9.5C2.375 10.1299 3.12567 10.7341 4.46186 11.1795C5.79806 11.6249 7.61033 11.8752 9.5 11.8752C11.3897 11.8752 13.2019 11.6249 14.5381 11.1795C15.8743 10.7341 16.625 10.1299 16.625 9.5"/></svg>
      </a>
      <a href="/reports" class="sidebar-icon {{if eq .CurrentPage "reports"}}active{{end}}" title="Reports">
        <svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M11.0832 1.5827H4.75047C4.33059 1.5827 3.9279 1.74953 3.631 2.04649C3.3341 2.34344 3.1673 2.7462 3.1673 3.16616V15.8338C3.1673 16.2538 3.3341 16.6566 3.631 16.9535C3.9279 17.2505 4.33059 17.4173 4.75047 17.4173H14.2495C14.6694 17.4173 15.0721 17.2505 15.369 16.9535C15.6659 16.6566 15.8327 16.2538 15.8327 15.8338V6.33308M11.0832 1.5827C11.3338 1.5823 11.5819 1.63147 11.8134 1.72741C12.0449 1.82334 12.2552 1.96413 12.432 2.14166L15.2723 4.98239C15.4502 5.15935 15.5914 5.36982 15.6876 5.60165C15.7838 5.83348 15.8331 6.08208 15.8327 6.33308M11.0832 1.5827V5.54135C11.0832 5.75133 11.1666 5.95271 11.315 6.10119C11.4635 6.24966 11.6648 6.33308 11.8748 6.33308L15.8327 6.33308M7.91682 7.12481H6.33365M12.6663 10.2917H6.33365M12.6663 13.4586H6.33365"/></svg>
      </a>
    </nav>
  </div>
  <a href="/settings" class="sidebar-icon {{if eq .CurrentPage "settings"}}active{{end}}" title="Settings">
    <svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M7.65586 3.2751C7.69948 2.81621 7.91262 2.39006 8.25365 2.07991C8.59467 1.76976 9.03908 1.5979 9.50005 1.5979C9.96101 1.5979 10.4054 1.76976 10.7464 2.07991C11.0875 2.39006 11.3006 2.81621 11.3442 3.2751C11.3705 3.57155 11.4677 3.85731 11.6278 4.10821C11.7878 4.3591 12.006 4.56774 12.2637 4.71647C12.5215 4.8652 12.8113 4.94964 13.1086 4.96264C13.4059 4.97563 13.702 4.91681 13.9718 4.79115C14.3907 4.60097 14.8653 4.57345 15.3034 4.71395C15.7414 4.85444 16.1115 5.1529 16.3416 5.55123C16.5718 5.94956 16.6455 6.41927 16.5484 6.86893C16.4513 7.3186 16.1903 7.71605 15.8164 7.98394C15.5728 8.15482 15.374 8.38184 15.2368 8.64579C15.0995 8.90975 15.0279 9.20287 15.0279 9.50038C15.0279 9.79788 15.0995 10.091 15.2368 10.355C15.374 10.6189 15.5728 10.8459 15.8164 11.0168C16.1903 11.2847 16.4513 11.6822 16.5484 12.1318C16.6455 12.5815 16.5718 13.0512 16.3416 13.4495C16.1115 13.8479 15.7414 14.1463 15.3034 14.2868C14.8653 14.4273 14.3907 14.3998 13.9718 14.2096C13.702 14.0839 13.4059 14.0251 13.1086 14.0381C12.8113 14.0511 12.5215 14.1355 12.2637 14.2843C12.006 14.433 11.7878 14.6416 11.6278 14.8925C11.4677 15.1434 11.3705 15.4292 11.3442 15.7256C11.3006 16.1845 11.0875 16.6107 10.7464 16.9208C10.4054 17.231 9.96101 17.4028 9.50005 17.4028C9.03908 17.4028 8.59467 17.231 8.25365 16.9208C7.91262 16.6107 7.69948 16.1845 7.65586 15.7256C7.62969 15.4291 7.53244 15.1432 7.37234 14.8922C7.21224 14.6413 6.994 14.4325 6.73613 14.2838C6.47825 14.1351 6.18832 14.0507 5.8909 14.0377C5.59348 14.0248 5.29733 14.0838 5.02753 14.2096C4.60865 14.3998 4.13399 14.4273 3.69594 14.2868C3.25789 14.1463 2.88779 13.8479 2.65766 13.4495C2.42753 13.0512 2.35384 12.5815 2.45094 12.1318C2.54803 11.6822 2.80896 11.2847 3.18294 11.0168C3.42648 10.8459 3.62527 10.6189 3.76251 10.355C3.89975 10.091 3.97141 9.79788 3.97141 9.50038C3.97141 9.20287 3.89975 8.90975 3.76251 8.64579C3.62527 8.38184 3.42648 8.15482 3.18294 7.98394C2.80949 7.71591 2.54902 7.31862 2.45216 6.86926C2.3553 6.4199 2.42897 5.95057 2.65885 5.5525C2.88873 5.15443 3.2584 4.85604 3.69601 4.71534C4.13363 4.57463 4.60793 4.60165 5.02673 4.79115C5.2965 4.91681 5.59257 4.97563 5.88989 4.96264C6.18721 4.94964 6.47701 4.8652 6.73478 4.71647C6.99256 4.56774 7.2107 4.3591 7.37075 4.10821C7.53081 3.85731 7.62806 3.57155 7.65428 3.2751"/></svg>
  </a>
</aside>

<!-- Main content -->
<main class="main">
  <div class="page-header">
    <div>
      <h1>Secrets Overview</h1>
      <div class="subtitle">Aggregate view across {{.ActiveProviders}} active providers</div>
    </div>
    <div class="header-search">
      <div class="search-box">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/></svg>
        <input id="search" type="text" placeholder="Search keys, environments…" oninput="filterRows()">
      </div>
    </div>
  </div>

  <!-- Stats -->
  <div class="stats-row">
    <a href="/" class="stat-card">
      <div class="stat-card-top"><span class="stat-card-label">Total Secrets</span></div>
      <div class="stat-n">{{index .Counts "total"}}</div>
    </a>
    <a href="/?status=warning" class="stat-card">
      <div class="stat-card-top">
        <span class="stat-card-label">Expiring Soon</span>
        <span class="stat-tag orange">Next 7 days</span>
      </div>
      <div class="stat-n">{{index .Counts "warning"}}</div>
    </a>
    <a href="/?status=overdue" class="stat-card">
      <div class="stat-card-top">
        <span class="stat-card-label">Overdue Rotation</span>
        <span class="stat-tag red">Action required</span>
      </div>
      <div class="stat-n">{{index .Counts "overdue"}}</div>
    </a>
    <div class="stat-card">
      <div class="stat-card-top">
        <span class="stat-card-label">Compliance Rate</span>
        <span class="stat-tag green">Target &gt;90%</span>
      </div>
      <div class="stat-n compliance">{{.ComplianceRate}}</div>
    </div>
  </div>

  <!-- Provider integrations -->
  <div class="section-header">Active Provider Integrations</div>
  <div class="provider-row">
    {{range .StoreTypes}}
    <div class="provider-card">
      <div class="provider-icon {{.}}">{{storeInitial .}}</div>
      <div>
        <div class="provider-name">{{storeTabLabel .}}</div>
        <div class="provider-sync"><span class="sync-dot"></span>Sync {{$.LastScan}}</div>
      </div>
    </div>
    {{end}}
  </div>

  <!-- Table -->
  <div class="table-section">
    <div class="table-top">
      <h2>Managed Secrets</h2>
      <div class="table-top-right">
        <select class="provider-select" onchange="navigateProvider(this.value)">
          <option value="">All providers</option>
          {{range .StorePairs}}
          <option value="store={{.StoreType}}&account={{.Account}}" {{if and (eq $.StoreFilter .StoreType) (eq $.AccountFilter .Account)}}selected{{end}}>
            {{storeTabLabel .StoreType}}{{if .Account}} · {{.Account}}{{end}} ({{.Count}})
          </option>
          {{end}}
        </select>
        {{if and (not .ShowExcluded) (gt .ExcludedCount 0)}}
        <a class="excl-link" href="?{{if .StoreFilter}}store={{.StoreFilter}}&account={{.AccountFilter}}&{{end}}{{if .StatusFilter}}status={{.StatusFilter}}&{{end}}excluded=1">{{.ExcludedCount}} excluded</a>
        {{else if .ShowExcluded}}
        <a class="excl-link" href="?{{if .StoreFilter}}store={{.StoreFilter}}&account={{.AccountFilter}}&{{end}}{{if .StatusFilter}}status={{.StatusFilter}}{{end}}">Hide excluded</a>
        {{end}}
      </div>
    </div>
    <table>
      <thead>
        <tr>
          <th>Secret Name</th>
          <th>Provider</th>
          <th>Last Rotated</th>
          <th>Policy</th>
          <th>Status</th>
          <th></th>
        </tr>
      </thead>
      <tbody id="secret-list">
      {{if .Secrets}}
        {{range .Secrets}}
        {{$short := shortName .Name}}
        <tr class="{{.Status}}{{if .Excluded}} excluded{{end}}" data-name="{{.Name}}" data-full-path="{{.FullPath}}">
          <td>
            <span class="td-name">{{$short}}</span>
            {{if ne $short .Name}}<span class="td-fullpath">{{.Name}}</span>{{end}}
          </td>
          <td class="td-provider">
            {{storeLabel .StoreType}}
            {{if .Account}}<div style="font-size:10px;color:var(--faint);font-family:var(--mono);margin-top:1px">{{.Account}}</div>{{end}}
          </td>
          <td class="td-date">{{if .LastRotated}}{{timeAgo .LastRotated}}{{else}}—{{end}}</td>
          <td class="td-policy">{{.MaxAgeDays}} days</td>
          <td class="td-status">
            {{if .Excluded}}
            <span class="status-badge excluded">excluded</span>
            {{else}}
            <span class="status-badge {{.Status}}">{{.Status}}</span>
            {{end}}
          </td>
          <td class="td-actions">
            <div class="row-menu">
              <button class="menu-btn" onclick="toggleMenu(this)">···</button>
              <div class="dropdown">
                {{$url := consoleURL .StoreType .FullPath .Region .Account}}
                {{if $url}}<a href="{{$url}}" target="_blank" rel="noopener" class="dd-item">Open in console ↗</a>{{end}}
                <a href="/history?name={{.Name}}" class="dd-item">View history</a>
                <button class="dd-item{{if .Excluded}} is-excluded{{end}}" onclick="toggleExclude(this)" data-excluded="{{if .Excluded}}1{{else}}0{{end}}">
                  {{if .Excluded}}Re-include{{else}}Exclude{{end}}
                </button>
              </div>
            </div>
          </td>
        </tr>
        {{end}}
      {{else}}
        <tr class="empty-row"><td colspan="6">No secrets found. Run <code>secretwatch scan</code> first.</td></tr>
      {{end}}
      </tbody>
      <tfoot>
        <tr><td colspan="6">
          <div class="tbl-footer">
            <span id="footer-count">Showing {{len .Secrets}} of {{index .Counts "total"}} secrets</span>
            <span></span>
          </div>
        </td></tr>
      </tfoot>
    </table>
  </div>
</main>

<!-- Right panel -->
<aside class="right-panel">
  <div class="rp-section">
    <div class="rp-title">Rotation Timeline</div>
    <div class="rp-subtitle">Recent rotation actions and audits</div>
    {{if .RecentRotations}}
    {{range .RecentRotations}}
    {{$short := shortName .Name}}
    <div class="timeline-item">
      <div class="tl-icon {{.StoreType}}">{{storeInitial .StoreType}}</div>
      <div class="tl-body">
        <div class="tl-name">{{$short}}</div>
        <div class="tl-action">Automated system rotation</div>
      </div>
      <div class="tl-time">{{timeAgo .LastRotated}}</div>
    </div>
    {{end}}
    {{else}}
    <div style="font-size:12px;color:var(--muted);padding:8px 0">No rotation history yet. Run <code>secretwatch scan</code>.</div>
    {{end}}
  </div>

  <div class="rp-divider"></div>

  <div class="rp-section">
    <div class="rp-title">Audit Compliance</div>
    <div class="compliance-rate">{{.ComplianceRate}}</div>
    <div class="compliance-label">of secrets within rotation policy</div>
    <div class="compliance-note">Secrets are considered compliant when rotated within their configured policy window.</div>
  </div>
</aside>

<script>
function navigateProvider(val) {
  const status = new URLSearchParams(window.location.search).get('status');
  let url = '/?' + (val || '');
  if (status) url += (val ? '&' : '') + 'status=' + status;
  window.location.href = url;
}
function filterRows() {
  const q = document.getElementById('search').value.toLowerCase();
  let n = 0;
  document.querySelectorAll('#secret-list tr[data-name]').forEach(r => {
    const show = r.dataset.name.toLowerCase().includes(q);
    r.style.display = show ? '' : 'none';
    if (show) n++;
  });
  const fc = document.getElementById('footer-count');
  if (fc) fc.textContent = 'Showing ' + n + ' secrets';
}
function toggleMenu(btn) {
  const dd = btn.nextElementSibling;
  const wasOpen = dd.classList.contains('open');
  document.querySelectorAll('.dropdown.open').forEach(d => d.classList.remove('open'));
  if (!wasOpen) dd.classList.add('open');
}
document.addEventListener('click', e => {
  if (!e.target.closest('.row-menu')) document.querySelectorAll('.dropdown.open').forEach(d => d.classList.remove('open'));
});
async function toggleExclude(btn) {
  const row = btn.closest('tr');
  const nowExcluded = btn.dataset.excluded !== '1';
  const res = await fetch('/api/exclude', { method: 'POST', body: new URLSearchParams({ full_path: row.dataset.fullPath, excluded: nowExcluded ? '1' : '0' }) });
  if (!res.ok) return;
  btn.dataset.excluded = nowExcluded ? '1' : '0';
  btn.textContent = nowExcluded ? 'Re-include' : 'Exclude';
  btn.classList.toggle('is-excluded', nowExcluded);
  row.classList.toggle('excluded', nowExcluded);
  btn.closest('.dropdown').classList.remove('open');
}
</script>
</body>
</html>
`

// ── Shared sidebar/CSS snippet used by all inner pages ─────────────────────
// (inlined per-template to avoid template composition complexity)

const sharedCSS = `
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
:root {
  --bg:#09090b; --sidebar-bg:#09090b; --card:#0c0c0e; --card2:#141416;
  --border:#1f1f23; --text:#ededf0; --sub:#a1a1aa; --muted:#71717a; --faint:#52525b;
  --green:#10b981; --green-bg:rgba(16,185,129,0.1);
  --orange:#f59e0b; --orange-bg:rgba(245,158,11,0.1);
  --red:#ef4444; --red-bg:rgba(239,68,68,0.1);
  --blue:#4285f4; --aws:#ff9900;
  --font:'Geist',system-ui,sans-serif; --mono:'Geist Mono',monospace;
}
html,body{height:100%;}
body{background:var(--bg);color:var(--text);font-family:var(--font);font-size:13px;line-height:1.5;-webkit-font-smoothing:antialiased;display:flex;height:100vh;overflow:hidden;}
.sidebar{width:64px;flex-shrink:0;background:var(--sidebar-bg);border-right:1px solid var(--border);display:flex;flex-direction:column;align-items:center;justify-content:space-between;padding:24px 0;}
.sidebar-top{display:flex;flex-direction:column;align-items:center;gap:32px;}
.sidebar-logo{width:36px;height:36px;background:var(--card);border:1px solid var(--border);border-radius:8px;display:flex;align-items:center;justify-content:center;font-size:18px;font-weight:700;color:var(--text);font-family:var(--mono);text-decoration:none;}
.sidebar-nav{display:flex;flex-direction:column;gap:12px;}
.sidebar-icon{width:48px;height:48px;border-radius:8px;border:1px solid transparent;display:flex;align-items:center;justify-content:center;color:var(--muted);cursor:pointer;transition:all .12s;text-decoration:none;}
.sidebar-icon:hover,.sidebar-icon.active{background:var(--card);border-color:var(--border);color:var(--text);}
.sidebar-icon svg{width:19px;height:19px;}
.main{flex:1;overflow-y:auto;display:flex;flex-direction:column;min-width:0;}
.page-header{display:flex;align-items:flex-start;justify-content:space-between;padding:24px 28px 20px;flex-shrink:0;}
.page-header h1{font-size:20px;font-weight:700;letter-spacing:-0.3px;}
.page-header .subtitle{font-size:13px;color:var(--sub);margin-top:2px;}
.status-badge{display:inline-flex;align-items:center;font-size:11px;font-weight:600;padding:2px 8px;border-radius:12px;}
.status-badge.ok{background:var(--green-bg);color:var(--green);}
.status-badge.warning{background:var(--orange-bg);color:var(--orange);}
.status-badge.overdue{background:var(--red-bg);color:var(--red);}
.status-badge.unknown,.status-badge.excluded{color:var(--muted);background:rgba(113,113,122,0.1);}
`

const sharedSidebarFn = `{{define "sidebar"}}
<aside class="sidebar">
  <div class="sidebar-top">
    <a href="/" class="sidebar-logo">S</a>
    <nav class="sidebar-nav">
      <a href="/" class="sidebar-icon {{if eq .CurrentPage "dashboard"}}active{{end}}" title="Overview"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M2.375 7.125H16.625M2.375 11.875H16.625M7.125 2.375V16.625M11.875 2.375V16.625M3.95833 2.375H15.0417C15.9161 2.375 16.625 3.08388 16.625 3.95833V15.0417C16.625 15.9161 15.9161 16.625 15.0417 16.625H3.95833C3.08388 16.625 2.375 15.9161 2.375 15.0417V3.95833C2.375 3.08388 3.08388 2.375 3.95833 2.375Z"/></svg></a>
      <a href="/secrets" class="sidebar-icon {{if eq .CurrentPage "secrets"}}active{{end}}" title="Secrets"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M11.0831 5.54166L13.4582 7.91685M7.44128 8.39158L2.04664 13.7864C1.74967 14.0833 1.58279 14.486 1.5827 14.9059V16.6256C1.5827 16.8356 1.66611 17.0369 1.81458 17.1854C1.96306 17.3339 2.16443 17.4173 2.3744 17.4173H4.7495C4.95947 17.4173 5.16085 17.3339 5.30932 17.1854C5.45779 17.0369 5.5412 16.8356 5.5412 16.6256V15.8338C5.5412 15.6239 5.62461 15.4225 5.77309 15.274C5.92156 15.1255 6.12293 15.0421 6.3329 15.0421H7.1246C7.33458 15.0421 7.53595 14.9587 7.68442 14.8102C7.83289 14.6617 7.91631 14.4604 7.91631 14.2504V13.4587C7.91631 13.2487 7.99972 13.0473 8.14819 12.8988C8.29666 12.7503 8.49803 12.6669 8.70801 12.6669H8.84418C9.26409 12.6668 9.66677 12.4999 9.96364 12.203L10.6081 11.5585M9.81652 2.13736C10.1825 1.79792 10.6632 1.6093 11.1624 1.6093C11.6616 1.6093 12.1423 1.79792 12.5083 2.13736L16.8627 6.49187C17.2021 6.85788 17.3907 7.33863 17.3907 7.83781C17.3907 8.337 17.2021 8.81775 16.8627 9.18376L13.9334 12.1132C13.5674 12.4526 13.0866 12.6412 12.5875 12.6412C12.0883 12.6412 11.6076 12.4526 11.2416 12.1132L6.88723 7.75864C6.5478 7.39263 6.35919 6.91189 6.35919 6.4127C6.35919 5.91352 6.5478 5.43277 6.88723 5.06676L9.81652 2.13736Z"/></svg></a>
      <a href="/providers" class="sidebar-icon {{if eq .CurrentPage "providers"}}active{{end}}" title="Providers"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M16.625 3.95789C16.625 5.26967 13.435 6.33308 9.5 6.33308C5.56497 6.33308 2.375 5.26967 2.375 3.95789M16.625 3.95789C16.625 2.64611 13.435 1.5827 9.5 1.5827C5.56497 1.5827 2.375 2.64611 2.375 3.95789M16.625 3.95789V15.0421C16.625 15.6721 15.8743 16.2762 14.5381 16.7216C13.2019 17.1671 11.3897 17.4173 9.5 17.4173C7.61033 17.4173 5.79806 17.1671 4.46186 16.7216C3.12567 16.2762 2.375 15.6721 2.375 15.0421V3.95789M2.375 9.5C2.375 10.1299 3.12567 10.7341 4.46186 11.1795C5.79806 11.6249 7.61033 11.8752 9.5 11.8752C11.3897 11.8752 13.2019 11.6249 14.5381 11.1795C15.8743 10.7341 16.625 10.1299 16.625 9.5"/></svg></a>
      <a href="/reports" class="sidebar-icon {{if eq .CurrentPage "reports"}}active{{end}}" title="Reports"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M11.0832 1.5827H4.75047C4.33059 1.5827 3.9279 1.74953 3.631 2.04649C3.3341 2.34344 3.1673 2.7462 3.1673 3.16616V15.8338C3.1673 16.2538 3.3341 16.6566 3.631 16.9535C3.9279 17.2505 4.33059 17.4173 4.75047 17.4173H14.2495C14.6694 17.4173 15.0721 17.2505 15.369 16.9535C15.6659 16.6566 15.8327 16.2538 15.8327 15.8338V6.33308M11.0832 1.5827C11.3338 1.5823 11.5819 1.63147 11.8134 1.72741C12.0449 1.82334 12.2552 1.96413 12.432 2.14166L15.2723 4.98239C15.4502 5.15935 15.5914 5.36982 15.6876 5.60165C15.7838 5.83348 15.8331 6.08208 15.8327 6.33308M11.0832 1.5827V5.54135C11.0832 5.75133 11.1666 5.95271 11.315 6.10119C11.4635 6.24966 11.6648 6.33308 11.8748 6.33308L15.8327 6.33308M7.91682 7.12481H6.33365M12.6663 10.2917H6.33365M12.6663 13.4586H6.33365"/></svg></a>
    </nav>
  </div>
  <a href="/settings" class="sidebar-icon {{if eq .CurrentPage "settings"}}active{{end}}" title="Settings"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M7.65586 3.2751C7.69948 2.81621 7.91262 2.39006 8.25365 2.07991C8.59467 1.76976 9.03908 1.5979 9.50005 1.5979C9.96101 1.5979 10.4054 1.76976 10.7464 2.07991C11.0875 2.39006 11.3006 2.81621 11.3442 3.2751C11.3705 3.57155 11.4677 3.85731 11.6278 4.10821C11.7878 4.3591 12.006 4.56774 12.2637 4.71647C12.5215 4.8652 12.8113 4.94964 13.1086 4.96264C13.4059 4.97563 13.702 4.91681 13.9718 4.79115C14.3907 4.60097 14.8653 4.57345 15.3034 4.71395C15.7414 4.85444 16.1115 5.1529 16.3416 5.55123C16.5718 5.94956 16.6455 6.41927 16.5484 6.86893C16.4513 7.3186 16.1903 7.71605 15.8164 7.98394C15.5728 8.15482 15.374 8.38184 15.2368 8.64579C15.0995 8.90975 15.0279 9.20287 15.0279 9.50038C15.0279 9.79788 15.0995 10.091 15.2368 10.355C15.374 10.6189 15.5728 10.8459 15.8164 11.0168C16.1903 11.2847 16.4513 11.6822 16.5484 12.1318C16.6455 12.5815 16.5718 13.0512 16.3416 13.4495C16.1115 13.8479 15.7414 14.1463 15.3034 14.2868C14.8653 14.4273 14.3907 14.3998 13.9718 14.2096C13.702 14.0839 13.4059 14.0251 13.1086 14.0381C12.8113 14.0511 12.5215 14.1355 12.2637 14.2843C12.006 14.433 11.7878 14.6416 11.6278 14.8925C11.4677 15.1434 11.3705 15.4292 11.3442 15.7256C11.3006 16.1845 11.0875 16.6107 10.7464 16.9208C10.4054 17.231 9.96101 17.4028 9.50005 17.4028C9.03908 17.4028 8.59467 17.231 8.25365 16.9208C7.91262 16.6107 7.69948 16.1845 7.65586 15.7256C7.62969 15.4291 7.53244 15.1432 7.37234 14.8922C7.21224 14.6413 6.994 14.4325 6.73613 14.2838C6.47825 14.1351 6.18832 14.0507 5.8909 14.0377C5.59348 14.0248 5.29733 14.0838 5.02753 14.2096C4.60865 14.3998 4.13399 14.4273 3.69594 14.2868C3.25789 14.1463 2.88779 13.8479 2.65766 13.4495C2.42753 13.0512 2.35384 12.5815 2.45094 12.1318C2.54803 11.6822 2.80896 11.2847 3.18294 11.0168C3.42648 10.8459 3.62527 10.6189 3.76251 10.355C3.89975 10.091 3.97141 9.79788 3.97141 9.50038C3.97141 9.20287 3.89975 8.90975 3.76251 8.64579C3.62527 8.38184 3.42648 8.15482 3.18294 7.98394C2.80949 7.71591 2.54902 7.31862 2.45216 6.86926C2.3553 6.4199 2.42897 5.95057 2.65885 5.5525C2.88873 5.15443 3.2584 4.85604 3.69601 4.71534C4.13363 4.57463 4.60793 4.60165 5.02673 4.79115C5.2965 4.91681 5.59257 4.97563 5.88989 4.96264C6.18721 4.94964 6.47701 4.8652 6.73478 4.71647C6.99256 4.56774 7.2107 4.3591 7.37075 4.10821C7.53081 3.85731 7.62806 3.57155 7.65428 3.2751"/></svg></a>
</aside>
{{end}}`

const secretsHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>SecretWatch — Secrets</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Geist:wght@400;500;600;700&family=Geist+Mono:wght@400;500&display=swap" rel="stylesheet">
<style>
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
:root {
  --bg:#09090b; --sidebar-bg:#09090b; --card:#0c0c0e; --card2:#141416;
  --border:#1f1f23; --text:#ededf0; --sub:#a1a1aa; --muted:#71717a; --faint:#52525b;
  --green:#10b981; --green-bg:rgba(16,185,129,0.1);
  --orange:#f59e0b; --orange-bg:rgba(245,158,11,0.1);
  --red:#ef4444; --red-bg:rgba(239,68,68,0.1);
  --blue:#4285f4; --aws:#ff9900;
  --font:'Geist',system-ui,sans-serif; --mono:'Geist Mono',monospace;
}
html,body{height:100%;}
body{background:var(--bg);color:var(--text);font-family:var(--font);font-size:13px;line-height:1.5;-webkit-font-smoothing:antialiased;display:flex;height:100vh;overflow:hidden;}
.sidebar{width:64px;flex-shrink:0;background:var(--sidebar-bg);border-right:1px solid var(--border);display:flex;flex-direction:column;align-items:center;justify-content:space-between;padding:24px 0;}
.sidebar-top{display:flex;flex-direction:column;align-items:center;gap:32px;}
.sidebar-logo{width:36px;height:36px;background:var(--card);border:1px solid var(--border);border-radius:8px;display:flex;align-items:center;justify-content:center;font-size:18px;font-weight:700;color:var(--text);font-family:var(--mono);text-decoration:none;}
.sidebar-nav{display:flex;flex-direction:column;gap:12px;}
.sidebar-icon{width:48px;height:48px;border-radius:8px;border:1px solid transparent;display:flex;align-items:center;justify-content:center;color:var(--muted);cursor:pointer;transition:all .12s;text-decoration:none;}
.sidebar-icon:hover,.sidebar-icon.active{background:var(--card);border-color:var(--border);color:var(--text);}
.sidebar-icon svg{width:19px;height:19px;}
.main{flex:1;overflow-y:auto;display:flex;flex-direction:column;min-width:0;}
.page-header{display:flex;align-items:flex-start;justify-content:space-between;padding:24px 28px 16px;}
.page-header h1{font-size:20px;font-weight:700;letter-spacing:-0.3px;}
.page-header .subtitle{font-size:13px;color:var(--sub);margin-top:2px;}
.toolbar{display:flex;align-items:center;gap:8px;padding:0 28px 16px;flex-wrap:wrap;}
.filter-tabs{display:flex;gap:2px;}
.ftab{font-size:12px;font-weight:500;color:var(--muted);padding:5px 12px;border-radius:6px;text-decoration:none;transition:all .12s;}
.ftab:hover{color:var(--text);background:var(--card);}
.ftab.active{color:var(--text);background:var(--card);}
.search-wrap{position:relative;}
.search-wrap svg{position:absolute;left:10px;top:50%;transform:translateY(-50%);color:var(--muted);pointer-events:none;}
#search{background:var(--card);border:1px solid var(--border);border-radius:7px;color:var(--text);font-family:var(--font);font-size:12px;padding:6px 12px 6px 30px;width:200px;outline:none;transition:border-color .15s;}
#search:focus{border-color:var(--blue);}
#search::placeholder{color:var(--muted);}
.tbl-wrap{margin:0 28px 28px;background:var(--card);border-radius:8px;overflow:hidden;}
table{width:100%;border-collapse:collapse;}
thead th{text-align:left;font-size:11px;font-weight:600;color:var(--muted);padding:10px 14px;border-bottom:1px solid var(--border);}
tbody tr{border-bottom:1px solid var(--border);transition:background .1s;}
tbody tr:last-child{border-bottom:none;}
tbody tr:hover{background:rgba(24,24,27,0.6);}
tbody tr.excluded td{opacity:.35;}
td{padding:10px 14px;vertical-align:middle;}
.td-name{font-family:var(--mono);font-size:13px;color:var(--text);display:block;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;max-width:260px;}
.td-fullpath{font-family:var(--mono);font-size:10px;color:var(--faint);display:block;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;max-width:260px;margin-top:1px;}
.td-type{font-size:12px;color:var(--sub);}
.td-date{font-size:12px;font-family:var(--mono);color:var(--sub);white-space:nowrap;}
.td-policy{font-size:12px;font-family:var(--mono);color:var(--muted);white-space:nowrap;}
.status-badge{display:inline-flex;align-items:center;font-size:11px;font-weight:600;padding:2px 8px;border-radius:12px;}
.status-badge.ok{background:var(--green-bg);color:var(--green);}
.status-badge.warning{background:var(--orange-bg);color:var(--orange);}
.status-badge.overdue{background:var(--red-bg);color:var(--red);}
.status-badge.unknown,.status-badge.excluded{color:var(--muted);background:rgba(113,113,122,0.1);}
.td-actions{text-align:right;width:44px;}
.row-menu{position:relative;display:inline-block;}
.menu-btn{background:none;border:none;cursor:pointer;color:var(--muted);font-size:16px;padding:4px 8px;border-radius:5px;font-family:var(--font);transition:all .1s;letter-spacing:1px;}
.menu-btn:hover{color:var(--text);background:var(--border);}
.dropdown{display:none;position:absolute;right:0;top:100%;margin-top:4px;background:#1a1a1d;border:1px solid var(--border);border-radius:8px;box-shadow:0 8px 24px rgba(0,0,0,.4);min-width:140px;z-index:100;overflow:hidden;}
.dropdown.open{display:block;}
.dd-item{display:block;width:100%;text-align:left;padding:9px 14px;font-size:13px;color:var(--sub);background:none;border:none;cursor:pointer;font-family:var(--font);text-decoration:none;transition:all .1s;}
.dd-item:hover{color:var(--text);background:rgba(255,255,255,0.04);}
.dd-item.active{color:var(--blue);}
.tbl-footer{padding:11px 14px;border-top:1px solid var(--border);font-size:12px;color:var(--muted);display:flex;align-items:center;justify-content:space-between;}
.excl-link{color:var(--muted);text-decoration:none;}
.excl-link:hover{color:var(--sub);}
.empty{text-align:center;padding:60px 24px;color:var(--sub);}
</style>
</head>
<body>
<aside class="sidebar">
  <div class="sidebar-top">
    <a href="/" class="sidebar-logo">S</a>
    <nav class="sidebar-nav">
      <a href="/" class="sidebar-icon" title="Overview"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M2.375 7.125H16.625M2.375 11.875H16.625M7.125 2.375V16.625M11.875 2.375V16.625M3.95833 2.375H15.0417C15.9161 2.375 16.625 3.08388 16.625 3.95833V15.0417C16.625 15.9161 15.9161 16.625 15.0417 16.625H3.95833C3.08388 16.625 2.375 15.9161 2.375 15.0417V3.95833C2.375 3.08388 3.08388 2.375 3.95833 2.375Z"/></svg></a>
      <a href="/secrets" class="sidebar-icon active" title="Secrets"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M11.0831 5.54166L13.4582 7.91685M7.44128 8.39158L2.04664 13.7864C1.74967 14.0833 1.58279 14.486 1.5827 14.9059V16.6256C1.5827 16.8356 1.66611 17.0369 1.81458 17.1854C1.96306 17.3339 2.16443 17.4173 2.3744 17.4173H4.7495C4.95947 17.4173 5.16085 17.3339 5.30932 17.1854C5.45779 17.0369 5.5412 16.8356 5.5412 16.6256V15.8338C5.5412 15.6239 5.62461 15.4225 5.77309 15.274C5.92156 15.1255 6.12293 15.0421 6.3329 15.0421H7.1246C7.33458 15.0421 7.53595 14.9587 7.68442 14.8102C7.83289 14.6617 7.91631 14.4604 7.91631 14.2504V13.4587C7.91631 13.2487 7.99972 13.0473 8.14819 12.8988C8.29666 12.7503 8.49803 12.6669 8.70801 12.6669H8.84418C9.26409 12.6668 9.66677 12.4999 9.96364 12.203L10.6081 11.5585M9.81652 2.13736C10.1825 1.79792 10.6632 1.6093 11.1624 1.6093C11.6616 1.6093 12.1423 1.79792 12.5083 2.13736L16.8627 6.49187C17.2021 6.85788 17.3907 7.33863 17.3907 7.83781C17.3907 8.337 17.2021 8.81775 16.8627 9.18376L13.9334 12.1132C13.5674 12.4526 13.0866 12.6412 12.5875 12.6412C12.0883 12.6412 11.6076 12.4526 11.2416 12.1132L6.88723 7.75864C6.5478 7.39263 6.35919 6.91189 6.35919 6.4127C6.35919 5.91352 6.5478 5.43277 6.88723 5.06676L9.81652 2.13736Z"/></svg></a>
      <a href="/providers" class="sidebar-icon" title="Providers"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M16.625 3.95789C16.625 5.26967 13.435 6.33308 9.5 6.33308C5.56497 6.33308 2.375 5.26967 2.375 3.95789M16.625 3.95789C16.625 2.64611 13.435 1.5827 9.5 1.5827C5.56497 1.5827 2.375 2.64611 2.375 3.95789M16.625 3.95789V15.0421C16.625 15.6721 15.8743 16.2762 14.5381 16.7216C13.2019 17.1671 11.3897 17.4173 9.5 17.4173C7.61033 17.4173 5.79806 17.1671 4.46186 16.7216C3.12567 16.2762 2.375 15.6721 2.375 15.0421V3.95789M2.375 9.5C2.375 10.1299 3.12567 10.7341 4.46186 11.1795C5.79806 11.6249 7.61033 11.8752 9.5 11.8752C11.3897 11.8752 13.2019 11.6249 14.5381 11.1795C15.8743 10.7341 16.625 10.1299 16.625 9.5"/></svg></a>
      <a href="/reports" class="sidebar-icon" title="Reports"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M11.0832 1.5827H4.75047C4.33059 1.5827 3.9279 1.74953 3.631 2.04649C3.3341 2.34344 3.1673 2.7462 3.1673 3.16616V15.8338C3.1673 16.2538 3.3341 16.6566 3.631 16.9535C3.9279 17.2505 4.33059 17.4173 4.75047 17.4173H14.2495C14.6694 17.4173 15.0721 17.2505 15.369 16.9535C15.6659 16.6566 15.8327 16.2538 15.8327 15.8338V6.33308M11.0832 1.5827C11.3338 1.5823 11.5819 1.63147 11.8134 1.72741C12.0449 1.82334 12.2552 1.96413 12.432 2.14166L15.2723 4.98239C15.4502 5.15935 15.5914 5.36982 15.6876 5.60165C15.7838 5.83348 15.8331 6.08208 15.8327 6.33308M11.0832 1.5827V5.54135C11.0832 5.75133 11.1666 5.95271 11.315 6.10119C11.4635 6.24966 11.6648 6.33308 11.8748 6.33308L15.8327 6.33308M7.91682 7.12481H6.33365M12.6663 10.2917H6.33365M12.6663 13.4586H6.33365"/></svg></a>
    </nav>
  </div>
  <a href="/settings" class="sidebar-icon" title="Settings"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M7.65586 3.2751C7.69948 2.81621 7.91262 2.39006 8.25365 2.07991C8.59467 1.76976 9.03908 1.5979 9.50005 1.5979C9.96101 1.5979 10.4054 1.76976 10.7464 2.07991C11.0875 2.39006 11.3006 2.81621 11.3442 3.2751C11.3705 3.57155 11.4677 3.85731 11.6278 4.10821C11.7878 4.3591 12.006 4.56774 12.2637 4.71647C12.5215 4.8652 12.8113 4.94964 13.1086 4.96264C13.4059 4.97563 13.702 4.91681 13.9718 4.79115C14.3907 4.60097 14.8653 4.57345 15.3034 4.71395C15.7414 4.85444 16.1115 5.1529 16.3416 5.55123C16.5718 5.94956 16.6455 6.41927 16.5484 6.86893C16.4513 7.3186 16.1903 7.71605 15.8164 7.98394C15.5728 8.15482 15.374 8.38184 15.2368 8.64579C15.0995 8.90975 15.0279 9.20287 15.0279 9.50038C15.0279 9.79788 15.0995 10.091 15.2368 10.355C15.374 10.6189 15.5728 10.8459 15.8164 11.0168C16.1903 11.2847 16.4513 11.6822 16.5484 12.1318C16.6455 12.5815 16.5718 13.0512 16.3416 13.4495C16.1115 13.8479 15.7414 14.1463 15.3034 14.2868C14.8653 14.4273 14.3907 14.3998 13.9718 14.2096C13.702 14.0839 13.4059 14.0251 13.1086 14.0381C12.8113 14.0511 12.5215 14.1355 12.2637 14.2843C12.006 14.433 11.7878 14.6416 11.6278 14.8925C11.4677 15.1434 11.3705 15.4292 11.3442 15.7256C11.3006 16.1845 11.0875 16.6107 10.7464 16.9208C10.4054 17.231 9.96101 17.4028 9.50005 17.4028C9.03908 17.4028 8.59467 17.231 8.25365 16.9208C7.91262 16.6107 7.69948 16.1845 7.65586 15.7256C7.62969 15.4291 7.53244 15.1432 7.37234 14.8922C7.21224 14.6413 6.994 14.4325 6.73613 14.2838C6.47825 14.1351 6.18832 14.0507 5.8909 14.0377C5.59348 14.0248 5.29733 14.0838 5.02753 14.2096C4.60865 14.3998 4.13399 14.4273 3.69594 14.2868C3.25789 14.1463 2.88779 13.8479 2.65766 13.4495C2.42753 13.0512 2.35384 12.5815 2.45094 12.1318C2.54803 11.6822 2.80896 11.2847 3.18294 11.0168C3.42648 10.8459 3.62527 10.6189 3.76251 10.355C3.89975 10.091 3.97141 9.79788 3.97141 9.50038C3.97141 9.20287 3.89975 8.90975 3.76251 8.64579C3.62527 8.38184 3.42648 8.15482 3.18294 7.98394C2.80949 7.71591 2.54902 7.31862 2.45216 6.86926C2.3553 6.4199 2.42897 5.95057 2.65885 5.5525C2.88873 5.15443 3.2584 4.85604 3.69601 4.71534C4.13363 4.57463 4.60793 4.60165 5.02673 4.79115C5.2965 4.91681 5.59257 4.97563 5.88989 4.96264C6.18721 4.94964 6.47701 4.8652 6.73478 4.71647C6.99256 4.56774 7.2107 4.3591 7.37075 4.10821C7.53081 3.85731 7.62806 3.57155 7.65428 3.2751"/></svg></a>
</aside>
<main class="main">
  <div class="page-header">
    <div>
      <h1>Secrets</h1>
      <div class="subtitle">{{index .Counts "total"}} secrets across {{len .StoreTypes}} stores · last scan {{.LastScan}}</div>
    </div>
    <div class="search-wrap">
      <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/></svg>
      <input id="search" type="text" placeholder="Filter secrets…" oninput="filterRows()">
    </div>
  </div>
  <div class="toolbar">
    <div class="filter-tabs">
      <a href="/secrets" class="ftab {{if eq .StatusFilter ""}}active{{end}}">All</a>
      <a href="/secrets?status=overdue" class="ftab {{if eq .StatusFilter "overdue"}}active{{end}}">Overdue</a>
      <a href="/secrets?status=warning" class="ftab {{if eq .StatusFilter "warning"}}active{{end}}">Expiring</a>
      <a href="/secrets?status=ok" class="ftab {{if eq .StatusFilter "ok"}}active{{end}}">Healthy</a>
      <a href="/secrets?status=unknown" class="ftab {{if eq .StatusFilter "unknown"}}active{{end}}">Unknown</a>
    </div>
  </div>
  <div class="tbl-wrap">
    <table>
      <thead><tr><th>Secret Name</th><th>Provider</th><th>Last Rotated</th><th>Policy</th><th>Status</th><th></th></tr></thead>
      <tbody id="secret-list">
      {{if .Secrets}}
        {{range .Secrets}}
        {{$short := shortName .Name}}
        <tr class="{{.Status}}{{if .Excluded}} excluded{{end}}" data-name="{{.Name}}" data-full-path="{{.FullPath}}">
          <td>
            {{if ne $short .Name}}<span class="td-name">{{$short}}</span><span class="td-fullpath">{{.Name}}</span>
            {{else}}<span class="td-name">{{.Name}}</span>{{end}}
          </td>
          <td class="td-type">
            {{storeLabel .StoreType}}
            {{if .Account}}<div style="font-size:10px;color:var(--faint);font-family:var(--mono);margin-top:1px">{{.Account}}</div>{{end}}
          </td>
          <td class="td-date">{{if .LastRotated}}{{timeAgo .LastRotated}}{{else}}—{{end}}</td>
          <td class="td-policy">{{.MaxAgeDays}} days</td>
          <td>{{if .Excluded}}<span class="status-badge excluded">excluded</span>{{else}}<span class="status-badge {{.Status}}">{{.Status}}</span>{{end}}</td>
          <td class="td-actions">
            <div class="row-menu">
              <button class="menu-btn" onclick="toggleMenu(this)">···</button>
              <div class="dropdown">
                {{$url := consoleURL .StoreType .FullPath .Region .Account}}
                {{if $url}}<a href="{{$url}}" target="_blank" rel="noopener" class="dd-item">Open in console ↗</a>{{end}}
                <a href="/history?name={{.Name}}" class="dd-item">View history</a>
                <button class="dd-item{{if .Excluded}} active{{end}}" onclick="toggleExclude(this)" data-excluded="{{if .Excluded}}1{{else}}0{{end}}">{{if .Excluded}}Re-include{{else}}Exclude{{end}}</button>
              </div>
            </div>
          </td>
        </tr>
        {{end}}
      {{else}}
        <tr><td colspan="6" class="empty">No secrets found. Run <code>secretwatch scan</code> first.</td></tr>
      {{end}}
      </tbody>
    </table>
    <div class="tbl-footer">
      <span id="footer-count">Showing {{len .Secrets}} secrets</span>
      {{if and (not .ShowExcluded) (gt .ExcludedCount 0)}}
      <a class="excl-link" href="/secrets?excluded=1{{if .StatusFilter}}&status={{.StatusFilter}}{{end}}">Show {{.ExcludedCount}} excluded</a>
      {{else if .ShowExcluded}}
      <a class="excl-link" href="/secrets{{if .StatusFilter}}?status={{.StatusFilter}}{{end}}">Hide excluded</a>
      {{end}}
    </div>
  </div>
</main>
<script>
function filterRows(){const q=document.getElementById('search').value.toLowerCase();let n=0;document.querySelectorAll('#secret-list tr[data-name]').forEach(r=>{const show=r.dataset.name.toLowerCase().includes(q);r.style.display=show?'':'none';if(show)n++;});document.getElementById('footer-count').textContent='Showing '+n+' secrets';}
function toggleMenu(btn){const dd=btn.nextElementSibling;const was=dd.classList.contains('open');document.querySelectorAll('.dropdown.open').forEach(d=>d.classList.remove('open'));if(!was)dd.classList.add('open');}
document.addEventListener('click',e=>{if(!e.target.closest('.row-menu'))document.querySelectorAll('.dropdown.open').forEach(d=>d.classList.remove('open'));});
async function toggleExclude(btn){const row=btn.closest('tr');const nowEx=btn.dataset.excluded!=='1';const res=await fetch('/api/exclude',{method:'POST',body:new URLSearchParams({full_path:row.dataset.fullPath,excluded:nowEx?'1':'0'})});if(!res.ok)return;btn.dataset.excluded=nowEx?'1':'0';btn.textContent=nowEx?'Re-include':'Exclude';btn.classList.toggle('active',nowEx);row.classList.toggle('excluded',nowEx);btn.closest('.dropdown').classList.remove('open');}
</script>
</body></html>
`

const providersHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>SecretWatch — Providers</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Geist:wght@400;500;600;700&family=Geist+Mono:wght@400;500&display=swap" rel="stylesheet">
<style>
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
:root {
  --bg:#09090b; --sidebar-bg:#09090b; --card:#0c0c0e; --card2:#141416;
  --border:#1f1f23; --text:#ededf0; --sub:#a1a1aa; --muted:#71717a; --faint:#52525b;
  --green:#10b981; --green-bg:rgba(16,185,129,0.1);
  --orange:#f59e0b; --orange-bg:rgba(245,158,11,0.1);
  --red:#ef4444; --red-bg:rgba(239,68,68,0.1);
  --blue:#4285f4; --aws:#ff9900;
  --font:'Geist',system-ui,sans-serif; --mono:'Geist Mono',monospace;
}
html,body{height:100%;}
body{background:var(--bg);color:var(--text);font-family:var(--font);font-size:13px;line-height:1.5;-webkit-font-smoothing:antialiased;display:flex;height:100vh;overflow:hidden;}
.sidebar{width:64px;flex-shrink:0;background:var(--sidebar-bg);border-right:1px solid var(--border);display:flex;flex-direction:column;align-items:center;justify-content:space-between;padding:24px 0;}
.sidebar-top{display:flex;flex-direction:column;align-items:center;gap:32px;}
.sidebar-logo{width:36px;height:36px;background:var(--card);border:1px solid var(--border);border-radius:8px;display:flex;align-items:center;justify-content:center;font-size:18px;font-weight:700;color:var(--text);font-family:var(--mono);text-decoration:none;}
.sidebar-nav{display:flex;flex-direction:column;gap:12px;}
.sidebar-icon{width:48px;height:48px;border-radius:8px;border:1px solid transparent;display:flex;align-items:center;justify-content:center;color:var(--muted);cursor:pointer;transition:all .12s;text-decoration:none;}
.sidebar-icon:hover,.sidebar-icon.active{background:var(--card);border-color:var(--border);color:var(--text);}
.sidebar-icon svg{width:19px;height:19px;}
.main{flex:1;overflow-y:auto;display:flex;flex-direction:column;min-width:0;}
.page-header{padding:24px 28px 24px;}
.page-header h1{font-size:20px;font-weight:700;letter-spacing:-0.3px;}
.page-header .subtitle{font-size:13px;color:var(--sub);margin-top:2px;}
.grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(280px,1fr));gap:16px;padding:0 28px 28px;}
.pcard{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:20px;}
.pcard-header{display:flex;align-items:center;gap:12px;margin-bottom:16px;}
.picon{width:36px;height:36px;border-radius:8px;display:flex;align-items:center;justify-content:center;font-size:13px;font-weight:700;font-family:var(--mono);flex-shrink:0;}
.picon.aws{background:rgba(255,153,0,0.15);color:var(--aws);}
.picon.gcp{background:rgba(66,133,244,0.15);color:var(--blue);}
.picon.vault{background:rgba(239,68,68,0.15);color:var(--red);}
.picon.github{background:rgba(161,161,170,0.15);color:var(--sub);}
.picon.other{background:rgba(113,113,122,0.15);color:var(--muted);}
.pname{font-size:14px;font-weight:600;}
.pcount{font-size:28px;font-weight:700;font-family:var(--mono);color:var(--text);margin-bottom:12px;}
.pcount-label{font-size:11px;color:var(--muted);margin-top:2px;}
.pstats{display:flex;gap:16px;padding-top:12px;border-top:1px solid var(--border);}
.pstat{text-align:center;}
.pstat-n{font-size:14px;font-weight:600;font-family:var(--mono);}
.pstat-n.red{color:var(--red);}
.pstat-n.orange{color:var(--orange);}
.pstat-n.green{color:var(--green);}
.pstat-n.muted{color:var(--muted);}
.pstat-l{font-size:11px;color:var(--muted);margin-top:1px;}
.empty-state{text-align:center;padding:80px 28px;color:var(--sub);}
</style>
</head>
<body>
<aside class="sidebar">
  <div class="sidebar-top">
    <a href="/" class="sidebar-logo">S</a>
    <nav class="sidebar-nav">
      <a href="/" class="sidebar-icon" title="Overview"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M2.375 7.125H16.625M2.375 11.875H16.625M7.125 2.375V16.625M11.875 2.375V16.625M3.95833 2.375H15.0417C15.9161 2.375 16.625 3.08388 16.625 3.95833V15.0417C16.625 15.9161 15.9161 16.625 15.0417 16.625H3.95833C3.08388 16.625 2.375 15.9161 2.375 15.0417V3.95833C2.375 3.08388 3.08388 2.375 3.95833 2.375Z"/></svg></a>
      <a href="/secrets" class="sidebar-icon" title="Secrets"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M11.0831 5.54166L13.4582 7.91685M7.44128 8.39158L2.04664 13.7864C1.74967 14.0833 1.58279 14.486 1.5827 14.9059V16.6256C1.5827 16.8356 1.66611 17.0369 1.81458 17.1854C1.96306 17.3339 2.16443 17.4173 2.3744 17.4173H4.7495C4.95947 17.4173 5.16085 17.3339 5.30932 17.1854C5.45779 17.0369 5.5412 16.8356 5.5412 16.6256V15.8338C5.5412 15.6239 5.62461 15.4225 5.77309 15.274C5.92156 15.1255 6.12293 15.0421 6.3329 15.0421H7.1246C7.33458 15.0421 7.53595 14.9587 7.68442 14.8102C7.83289 14.6617 7.91631 14.4604 7.91631 14.2504V13.4587C7.91631 13.2487 7.99972 13.0473 8.14819 12.8988C8.29666 12.7503 8.49803 12.6669 8.70801 12.6669H8.84418C9.26409 12.6668 9.66677 12.4999 9.96364 12.203L10.6081 11.5585M9.81652 2.13736C10.1825 1.79792 10.6632 1.6093 11.1624 1.6093C11.6616 1.6093 12.1423 1.79792 12.5083 2.13736L16.8627 6.49187C17.2021 6.85788 17.3907 7.33863 17.3907 7.83781C17.3907 8.337 17.2021 8.81775 16.8627 9.18376L13.9334 12.1132C13.5674 12.4526 13.0866 12.6412 12.5875 12.6412C12.0883 12.6412 11.6076 12.4526 11.2416 12.1132L6.88723 7.75864C6.5478 7.39263 6.35919 6.91189 6.35919 6.4127C6.35919 5.91352 6.5478 5.43277 6.88723 5.06676L9.81652 2.13736Z"/></svg></a>
      <a href="/providers" class="sidebar-icon active" title="Providers"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M16.625 3.95789C16.625 5.26967 13.435 6.33308 9.5 6.33308C5.56497 6.33308 2.375 5.26967 2.375 3.95789M16.625 3.95789C16.625 2.64611 13.435 1.5827 9.5 1.5827C5.56497 1.5827 2.375 2.64611 2.375 3.95789M16.625 3.95789V15.0421C16.625 15.6721 15.8743 16.2762 14.5381 16.7216C13.2019 17.1671 11.3897 17.4173 9.5 17.4173C7.61033 17.4173 5.79806 17.1671 4.46186 16.7216C3.12567 16.2762 2.375 15.6721 2.375 15.0421V3.95789M2.375 9.5C2.375 10.1299 3.12567 10.7341 4.46186 11.1795C5.79806 11.6249 7.61033 11.8752 9.5 11.8752C11.3897 11.8752 13.2019 11.6249 14.5381 11.1795C15.8743 10.7341 16.625 10.1299 16.625 9.5"/></svg></a>
      <a href="/reports" class="sidebar-icon" title="Reports"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M11.0832 1.5827H4.75047C4.33059 1.5827 3.9279 1.74953 3.631 2.04649C3.3341 2.34344 3.1673 2.7462 3.1673 3.16616V15.8338C3.1673 16.2538 3.3341 16.6566 3.631 16.9535C3.9279 17.2505 4.33059 17.4173 4.75047 17.4173H14.2495C14.6694 17.4173 15.0721 17.2505 15.369 16.9535C15.6659 16.6566 15.8327 16.2538 15.8327 15.8338V6.33308M11.0832 1.5827C11.3338 1.5823 11.5819 1.63147 11.8134 1.72741C12.0449 1.82334 12.2552 1.96413 12.432 2.14166L15.2723 4.98239C15.4502 5.15935 15.5914 5.36982 15.6876 5.60165C15.7838 5.83348 15.8331 6.08208 15.8327 6.33308M11.0832 1.5827V5.54135C11.0832 5.75133 11.1666 5.95271 11.315 6.10119C11.4635 6.24966 11.6648 6.33308 11.8748 6.33308L15.8327 6.33308M7.91682 7.12481H6.33365M12.6663 10.2917H6.33365M12.6663 13.4586H6.33365"/></svg></a>
    </nav>
  </div>
  <a href="/settings" class="sidebar-icon" title="Settings"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M7.65586 3.2751C7.69948 2.81621 7.91262 2.39006 8.25365 2.07991C8.59467 1.76976 9.03908 1.5979 9.50005 1.5979C9.96101 1.5979 10.4054 1.76976 10.7464 2.07991C11.0875 2.39006 11.3006 2.81621 11.3442 3.2751C11.3705 3.57155 11.4677 3.85731 11.6278 4.10821C11.7878 4.3591 12.006 4.56774 12.2637 4.71647C12.5215 4.8652 12.8113 4.94964 13.1086 4.96264C13.4059 4.97563 13.702 4.91681 13.9718 4.79115C14.3907 4.60097 14.8653 4.57345 15.3034 4.71395C15.7414 4.85444 16.1115 5.1529 16.3416 5.55123C16.5718 5.94956 16.6455 6.41927 16.5484 6.86893C16.4513 7.3186 16.1903 7.71605 15.8164 7.98394C15.5728 8.15482 15.374 8.38184 15.2368 8.64579C15.0995 8.90975 15.0279 9.20287 15.0279 9.50038C15.0279 9.79788 15.0995 10.091 15.2368 10.355C15.374 10.6189 15.5728 10.8459 15.8164 11.0168C16.1903 11.2847 16.4513 11.6822 16.5484 12.1318C16.6455 12.5815 16.5718 13.0512 16.3416 13.4495C16.1115 13.8479 15.7414 14.1463 15.3034 14.2868C14.8653 14.4273 14.3907 14.3998 13.9718 14.2096C13.702 14.0839 13.4059 14.0251 13.1086 14.0381C12.8113 14.0511 12.5215 14.1355 12.2637 14.2843C12.006 14.433 11.7878 14.6416 11.6278 14.8925C11.4677 15.1434 11.3705 15.4292 11.3442 15.7256C11.3006 16.1845 11.0875 16.6107 10.7464 16.9208C10.4054 17.231 9.96101 17.4028 9.50005 17.4028C9.03908 17.4028 8.59467 17.231 8.25365 16.9208C7.91262 16.6107 7.69948 16.1845 7.65586 15.7256C7.62969 15.4291 7.53244 15.1432 7.37234 14.8922C7.21224 14.6413 6.994 14.4325 6.73613 14.2838C6.47825 14.1351 6.18832 14.0507 5.8909 14.0377C5.59348 14.0248 5.29733 14.0838 5.02753 14.2096C4.60865 14.3998 4.13399 14.4273 3.69594 14.2868C3.25789 14.1463 2.88779 13.8479 2.65766 13.4495C2.42753 13.0512 2.35384 12.5815 2.45094 12.1318C2.54803 11.6822 2.80896 11.2847 3.18294 11.0168C3.42648 10.8459 3.62527 10.6189 3.76251 10.355C3.89975 10.091 3.97141 9.79788 3.97141 9.50038C3.97141 9.20287 3.89975 8.90975 3.76251 8.64579C3.62527 8.38184 3.42648 8.15482 3.18294 7.98394C2.80949 7.71591 2.54902 7.31862 2.45216 6.86926C2.3553 6.4199 2.42897 5.95057 2.65885 5.5525C2.88873 5.15443 3.2584 4.85604 3.69601 4.71534C4.13363 4.57463 4.60793 4.60165 5.02673 4.79115C5.2965 4.91681 5.59257 4.97563 5.88989 4.96264C6.18721 4.94964 6.47701 4.8652 6.73478 4.71647C6.99256 4.56774 7.2107 4.3591 7.37075 4.10821C7.53081 3.85731 7.62806 3.57155 7.65428 3.2751"/></svg></a>
</aside>
<main class="main">
  <div class="page-header">
    <h1>Providers</h1>
    <div class="subtitle">Active store integrations · last scan {{.LastScan}}</div>
  </div>
  {{if .Stats}}
  <div class="grid">
    {{range .Stats}}
    <div class="pcard">
      <div class="pcard-header">
        <div class="picon {{.StoreType}}">{{storeInitial .StoreType}}</div>
        <div class="pname">{{storeTabLabel .StoreType}}</div>
      </div>
      <div class="pcount">{{.Total}}<div class="pcount-label">secrets tracked</div></div>
      <div class="pstats">
        <div class="pstat"><div class="pstat-n red">{{.Overdue}}</div><div class="pstat-l">overdue</div></div>
        <div class="pstat"><div class="pstat-n orange">{{.Warning}}</div><div class="pstat-l">expiring</div></div>
        <div class="pstat"><div class="pstat-n green">{{.OK}}</div><div class="pstat-l">healthy</div></div>
        <div class="pstat"><div class="pstat-n muted">{{.Unknown}}</div><div class="pstat-l">unknown</div></div>
      </div>
    </div>
    {{end}}
  </div>
  {{else}}
  <div class="empty-state">No providers found. Run <code>secretwatch scan</code> first.</div>
  {{end}}
</main>
</body></html>
`

const reportsHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>SecretWatch — Reports</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Geist:wght@400;500;600;700&family=Geist+Mono:wght@400;500&display=swap" rel="stylesheet">
<style>
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
:root {
  --bg:#09090b; --sidebar-bg:#09090b; --card:#0c0c0e; --card2:#141416;
  --border:#1f1f23; --text:#ededf0; --sub:#a1a1aa; --muted:#71717a; --faint:#52525b;
  --green:#10b981; --green-bg:rgba(16,185,129,0.1);
  --orange:#f59e0b; --orange-bg:rgba(245,158,11,0.1);
  --red:#ef4444; --red-bg:rgba(239,68,68,0.1);
  --blue:#4285f4; --aws:#ff9900;
  --font:'Geist',system-ui,sans-serif; --mono:'Geist Mono',monospace;
}
html,body{height:100%;}
body{background:var(--bg);color:var(--text);font-family:var(--font);font-size:13px;line-height:1.5;-webkit-font-smoothing:antialiased;display:flex;height:100vh;overflow:hidden;}
.sidebar{width:64px;flex-shrink:0;background:var(--sidebar-bg);border-right:1px solid var(--border);display:flex;flex-direction:column;align-items:center;justify-content:space-between;padding:24px 0;}
.sidebar-top{display:flex;flex-direction:column;align-items:center;gap:32px;}
.sidebar-logo{width:36px;height:36px;background:var(--card);border:1px solid var(--border);border-radius:8px;display:flex;align-items:center;justify-content:center;font-size:18px;font-weight:700;color:var(--text);font-family:var(--mono);text-decoration:none;}
.sidebar-nav{display:flex;flex-direction:column;gap:12px;}
.sidebar-icon{width:48px;height:48px;border-radius:8px;border:1px solid transparent;display:flex;align-items:center;justify-content:center;color:var(--muted);cursor:pointer;transition:all .12s;text-decoration:none;}
.sidebar-icon:hover,.sidebar-icon.active{background:var(--card);border-color:var(--border);color:var(--text);}
.sidebar-icon svg{width:19px;height:19px;}
.main{flex:1;overflow-y:auto;display:flex;flex-direction:column;min-width:0;}
.page-header{padding:24px 28px 24px;}
.page-header h1{font-size:20px;font-weight:700;letter-spacing:-0.3px;}
.page-header .subtitle{font-size:13px;color:var(--sub);margin-top:2px;}
.content{padding:0 28px 28px;display:grid;grid-template-columns:1fr 1fr;gap:20px;align-items:start;}
.section{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:24px;}
.section h2{font-size:14px;font-weight:600;margin-bottom:4px;}
.section-sub{font-size:12px;color:var(--muted);margin-bottom:20px;}
.compliance-big{font-size:48px;font-weight:700;font-family:var(--mono);color:var(--green);line-height:1;margin-bottom:8px;}
.compliance-meta{font-size:12px;color:var(--muted);}
.stat-row{display:flex;justify-content:space-between;padding:8px 0;border-bottom:1px solid var(--border);font-size:13px;}
.stat-row:last-child{border-bottom:none;}
.stat-row-label{color:var(--sub);}
.stat-row-val{font-weight:600;font-family:var(--mono);}
.stat-row-val.red{color:var(--red);}
.stat-row-val.orange{color:var(--orange);}
.stat-row-val.green{color:var(--green);}
table{width:100%;border-collapse:collapse;margin-top:16px;}
thead th{text-align:left;font-size:11px;font-weight:600;color:var(--muted);padding:8px 0;border-bottom:1px solid var(--border);}
tbody tr{border-bottom:1px solid var(--border);}
tbody tr:last-child{border-bottom:none;}
td{padding:10px 0;font-size:12px;color:var(--sub);}
td:first-child{color:var(--text);font-weight:500;}
td.mono{font-family:var(--mono);}
td.green{color:var(--green);}
td.red{color:var(--red);}
.export-links{display:flex;flex-direction:column;gap:10px;margin-top:16px;}
.export-btn{display:block;padding:10px 14px;background:var(--card2);border:1px solid var(--border);border-radius:7px;color:var(--text);text-decoration:none;font-size:13px;transition:border-color .12s;}
.export-btn:hover{border-color:var(--sub);}
.export-btn span{display:block;font-size:11px;color:var(--muted);margin-top:2px;}
.cmd{font-family:var(--mono);font-size:12px;background:var(--card2);border:1px solid var(--border);border-radius:6px;padding:10px 14px;color:var(--sub);margin-top:8px;}
</style>
</head>
<body>
<aside class="sidebar">
  <div class="sidebar-top">
    <a href="/" class="sidebar-logo">S</a>
    <nav class="sidebar-nav">
      <a href="/" class="sidebar-icon" title="Overview"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M2.375 7.125H16.625M2.375 11.875H16.625M7.125 2.375V16.625M11.875 2.375V16.625M3.95833 2.375H15.0417C15.9161 2.375 16.625 3.08388 16.625 3.95833V15.0417C16.625 15.9161 15.9161 16.625 15.0417 16.625H3.95833C3.08388 16.625 2.375 15.9161 2.375 15.0417V3.95833C2.375 3.08388 3.08388 2.375 3.95833 2.375Z"/></svg></a>
      <a href="/secrets" class="sidebar-icon" title="Secrets"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M11.0831 5.54166L13.4582 7.91685M7.44128 8.39158L2.04664 13.7864C1.74967 14.0833 1.58279 14.486 1.5827 14.9059V16.6256C1.5827 16.8356 1.66611 17.0369 1.81458 17.1854C1.96306 17.3339 2.16443 17.4173 2.3744 17.4173H4.7495C4.95947 17.4173 5.16085 17.3339 5.30932 17.1854C5.45779 17.0369 5.5412 16.8356 5.5412 16.6256V15.8338C5.5412 15.6239 5.62461 15.4225 5.77309 15.274C5.92156 15.1255 6.12293 15.0421 6.3329 15.0421H7.1246C7.33458 15.0421 7.53595 14.9587 7.68442 14.8102C7.83289 14.6617 7.91631 14.4604 7.91631 14.2504V13.4587C7.91631 13.2487 7.99972 13.0473 8.14819 12.8988C8.29666 12.7503 8.49803 12.6669 8.70801 12.6669H8.84418C9.26409 12.6668 9.66677 12.4999 9.96364 12.203L10.6081 11.5585M9.81652 2.13736C10.1825 1.79792 10.6632 1.6093 11.1624 1.6093C11.6616 1.6093 12.1423 1.79792 12.5083 2.13736L16.8627 6.49187C17.2021 6.85788 17.3907 7.33863 17.3907 7.83781C17.3907 8.337 17.2021 8.81775 16.8627 9.18376L13.9334 12.1132C13.5674 12.4526 13.0866 12.6412 12.5875 12.6412C12.0883 12.6412 11.6076 12.4526 11.2416 12.1132L6.88723 7.75864C6.5478 7.39263 6.35919 6.91189 6.35919 6.4127C6.35919 5.91352 6.5478 5.43277 6.88723 5.06676L9.81652 2.13736Z"/></svg></a>
      <a href="/providers" class="sidebar-icon" title="Providers"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M16.625 3.95789C16.625 5.26967 13.435 6.33308 9.5 6.33308C5.56497 6.33308 2.375 5.26967 2.375 3.95789M16.625 3.95789C16.625 2.64611 13.435 1.5827 9.5 1.5827C5.56497 1.5827 2.375 2.64611 2.375 3.95789M16.625 3.95789V15.0421C16.625 15.6721 15.8743 16.2762 14.5381 16.7216C13.2019 17.1671 11.3897 17.4173 9.5 17.4173C7.61033 17.4173 5.79806 17.1671 4.46186 16.7216C3.12567 16.2762 2.375 15.6721 2.375 15.0421V3.95789M2.375 9.5C2.375 10.1299 3.12567 10.7341 4.46186 11.1795C5.79806 11.6249 7.61033 11.8752 9.5 11.8752C11.3897 11.8752 13.2019 11.6249 14.5381 11.1795C15.8743 10.7341 16.625 10.1299 16.625 9.5"/></svg></a>
      <a href="/reports" class="sidebar-icon active" title="Reports"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M11.0832 1.5827H4.75047C4.33059 1.5827 3.9279 1.74953 3.631 2.04649C3.3341 2.34344 3.1673 2.7462 3.1673 3.16616V15.8338C3.1673 16.2538 3.3341 16.6566 3.631 16.9535C3.9279 17.2505 4.33059 17.4173 4.75047 17.4173H14.2495C14.6694 17.4173 15.0721 17.2505 15.369 16.9535C15.6659 16.6566 15.8327 16.2538 15.8327 15.8338V6.33308M11.0832 1.5827C11.3338 1.5823 11.5819 1.63147 11.8134 1.72741C12.0449 1.82334 12.2552 1.96413 12.432 2.14166L15.2723 4.98239C15.4502 5.15935 15.5914 5.36982 15.6876 5.60165C15.7838 5.83348 15.8331 6.08208 15.8327 6.33308M11.0832 1.5827V5.54135C11.0832 5.75133 11.1666 5.95271 11.315 6.10119C11.4635 6.24966 11.6648 6.33308 11.8748 6.33308L15.8327 6.33308M7.91682 7.12481H6.33365M12.6663 10.2917H6.33365M12.6663 13.4586H6.33365"/></svg></a>
    </nav>
  </div>
  <a href="/settings" class="sidebar-icon" title="Settings"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M7.65586 3.2751C7.69948 2.81621 7.91262 2.39006 8.25365 2.07991C8.59467 1.76976 9.03908 1.5979 9.50005 1.5979C9.96101 1.5979 10.4054 1.76976 10.7464 2.07991C11.0875 2.39006 11.3006 2.81621 11.3442 3.2751C11.3705 3.57155 11.4677 3.85731 11.6278 4.10821C11.7878 4.3591 12.006 4.56774 12.2637 4.71647C12.5215 4.8652 12.8113 4.94964 13.1086 4.96264C13.4059 4.97563 13.702 4.91681 13.9718 4.79115C14.3907 4.60097 14.8653 4.57345 15.3034 4.71395C15.7414 4.85444 16.1115 5.1529 16.3416 5.55123C16.5718 5.94956 16.6455 6.41927 16.5484 6.86893C16.4513 7.3186 16.1903 7.71605 15.8164 7.98394C15.5728 8.15482 15.374 8.38184 15.2368 8.64579C15.0995 8.90975 15.0279 9.20287 15.0279 9.50038C15.0279 9.79788 15.0995 10.091 15.2368 10.355C15.374 10.6189 15.5728 10.8459 15.8164 11.0168C16.1903 11.2847 16.4513 11.6822 16.5484 12.1318C16.6455 12.5815 16.5718 13.0512 16.3416 13.4495C16.1115 13.8479 15.7414 14.1463 15.3034 14.2868C14.8653 14.4273 14.3907 14.3998 13.9718 14.2096C13.702 14.0839 13.4059 14.0251 13.1086 14.0381C12.8113 14.0511 12.5215 14.1355 12.2637 14.2843C12.006 14.433 11.7878 14.6416 11.6278 14.8925C11.4677 15.1434 11.3705 15.4292 11.3442 15.7256C11.3006 16.1845 11.0875 16.6107 10.7464 16.9208C10.4054 17.231 9.96101 17.4028 9.50005 17.4028C9.03908 17.4028 8.59467 17.231 8.25365 16.9208C7.91262 16.6107 7.69948 16.1845 7.65586 15.7256C7.62969 15.4291 7.53244 15.1432 7.37234 14.8922C7.21224 14.6413 6.994 14.4325 6.73613 14.2838C6.47825 14.1351 6.18832 14.0507 5.8909 14.0377C5.59348 14.0248 5.29733 14.0838 5.02753 14.2096C4.60865 14.3998 4.13399 14.4273 3.69594 14.2868C3.25789 14.1463 2.88779 13.8479 2.65766 13.4495C2.42753 13.0512 2.35384 12.5815 2.45094 12.1318C2.54803 11.6822 2.80896 11.2847 3.18294 11.0168C3.42648 10.8459 3.62527 10.6189 3.76251 10.355C3.89975 10.091 3.97141 9.79788 3.97141 9.50038C3.97141 9.20287 3.89975 8.90975 3.76251 8.64579C3.62527 8.38184 3.42648 8.15482 3.18294 7.98394C2.80949 7.71591 2.54902 7.31862 2.45216 6.86926C2.3553 6.4199 2.42897 5.95057 2.65885 5.5525C2.88873 5.15443 3.2584 4.85604 3.69601 4.71534C4.13363 4.57463 4.60793 4.60165 5.02673 4.79115C5.2965 4.91681 5.59257 4.97563 5.88989 4.96264C6.18721 4.94964 6.47701 4.8652 6.73478 4.71647C6.99256 4.56774 7.2107 4.3591 7.37075 4.10821C7.53081 3.85731 7.62806 3.57155 7.65428 3.2751"/></svg></a>
</aside>
<main class="main">
  <div class="page-header">
    <h1>Reports</h1>
    <div class="subtitle">Compliance summary and export · last scan {{.LastScan}}</div>
  </div>
  <div class="content">
    <div class="section">
      <h2>Compliance Rate</h2>
      <div class="section-sub">Healthy secrets as a percentage of non-unknown secrets</div>
      <div class="compliance-big">{{.ComplianceRate}}</div>
      <div class="compliance-meta">Target: &gt;90%</div>
      <div style="margin-top:20px">
        <div class="stat-row"><span class="stat-row-label">Total secrets</span><span class="stat-row-val">{{index .Counts "total"}}</span></div>
        <div class="stat-row"><span class="stat-row-label">Overdue</span><span class="stat-row-val red">{{index .Counts "overdue"}}</span></div>
        <div class="stat-row"><span class="stat-row-label">Expiring soon</span><span class="stat-row-val orange">{{index .Counts "warning"}}</span></div>
        <div class="stat-row"><span class="stat-row-label">Healthy</span><span class="stat-row-val green">{{index .Counts "ok"}}</span></div>
        <div class="stat-row"><span class="stat-row-label">Unknown</span><span class="stat-row-val">{{index .Counts "unknown"}}</span></div>
      </div>
      {{if .Stats}}
      <table>
        <thead><tr><th>Provider</th><th>Total</th><th>Healthy</th><th>Overdue</th></tr></thead>
        <tbody>
        {{range .Stats}}
        <tr>
          <td>{{storeTabLabel .StoreType}}</td>
          <td class="mono">{{.Total}}</td>
          <td class="mono green">{{.OK}}</td>
          <td class="mono{{if gt .Overdue 0}} red{{end}}">{{.Overdue}}</td>
        </tr>
        {{end}}
        </tbody>
      </table>
      {{end}}
    </div>
    <div class="section">
      <h2>Export</h2>
      <div class="section-sub">Download current scan data for auditing or integration</div>
      <div class="export-links">
        <a href="/api/secrets" class="export-btn" target="_blank">
          Download JSON
          <span>All secrets as JSON via /api/secrets</span>
        </a>
        <a href="/api/secrets?status=overdue" class="export-btn" target="_blank">
          Overdue only (JSON)
          <span>Filtered to overdue secrets only</span>
        </a>
      </div>
      <div style="margin-top:24px">
        <div style="font-size:12px;color:var(--muted);margin-bottom:8px">CLI equivalents</div>
        <div class="cmd">secretwatch list</div>
        <div class="cmd" style="margin-top:6px">secretwatch list --status overdue</div>
        <div class="cmd" style="margin-top:6px">secretwatch alert --dry-run</div>
      </div>
    </div>
  </div>
</main>
</body></html>
`

const settingsHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>SecretWatch — Settings</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Geist:wght@400;500;600;700&family=Geist+Mono:wght@400;500&display=swap" rel="stylesheet">
<style>
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
:root {
  --bg:#09090b; --sidebar-bg:#09090b; --card:#0c0c0e; --card2:#141416;
  --border:#1f1f23; --text:#ededf0; --sub:#a1a1aa; --muted:#71717a; --faint:#52525b;
  --green:#10b981; --font:'Geist',system-ui,sans-serif; --mono:'Geist Mono',monospace;
}
html,body{height:100%;}
body{background:var(--bg);color:var(--text);font-family:var(--font);font-size:13px;line-height:1.5;-webkit-font-smoothing:antialiased;display:flex;height:100vh;overflow:hidden;}
.sidebar{width:64px;flex-shrink:0;background:var(--sidebar-bg);border-right:1px solid var(--border);display:flex;flex-direction:column;align-items:center;justify-content:space-between;padding:24px 0;}
.sidebar-top{display:flex;flex-direction:column;align-items:center;gap:32px;}
.sidebar-logo{width:36px;height:36px;background:var(--card);border:1px solid var(--border);border-radius:8px;display:flex;align-items:center;justify-content:center;font-size:18px;font-weight:700;color:var(--text);font-family:var(--mono);text-decoration:none;}
.sidebar-nav{display:flex;flex-direction:column;gap:12px;}
.sidebar-icon{width:48px;height:48px;border-radius:8px;border:1px solid transparent;display:flex;align-items:center;justify-content:center;color:var(--muted);cursor:pointer;transition:all .12s;text-decoration:none;}
.sidebar-icon:hover,.sidebar-icon.active{background:var(--card);border-color:var(--border);color:var(--text);}
.sidebar-icon svg{width:19px;height:19px;}
.main{flex:1;overflow-y:auto;display:flex;flex-direction:column;min-width:0;}
.page-header{padding:24px 28px 24px;}
.page-header h1{font-size:20px;font-weight:700;letter-spacing:-0.3px;}
.page-header .subtitle{font-size:13px;color:var(--sub);margin-top:2px;}
.content{padding:0 28px 28px;display:flex;flex-direction:column;gap:16px;max-width:700px;}
.section{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:20px;}
.section h2{font-size:13px;font-weight:600;margin-bottom:16px;color:var(--sub);text-transform:uppercase;letter-spacing:0.5px;font-size:11px;}
.row{display:flex;justify-content:space-between;align-items:baseline;padding:7px 0;border-bottom:1px solid var(--border);font-size:13px;}
.row:last-child{border-bottom:none;padding-bottom:0;}
.row-label{color:var(--sub);}
.row-val{font-family:var(--mono);font-size:12px;color:var(--text);}
.row-val.muted{color:var(--muted);}
.notice{font-size:12px;color:var(--muted);background:var(--card2);border:1px solid var(--border);border-radius:7px;padding:12px 14px;line-height:1.6;}
.notice code{font-family:var(--mono);color:var(--sub);}
.badge-on{display:inline-block;font-size:10px;font-weight:600;padding:1px 7px;border-radius:10px;background:rgba(16,185,129,0.1);color:var(--green);}
.badge-off{display:inline-block;font-size:10px;font-weight:600;padding:1px 7px;border-radius:10px;background:rgba(113,113,122,0.1);color:var(--muted);}
</style>
</head>
<body>
<aside class="sidebar">
  <div class="sidebar-top">
    <a href="/" class="sidebar-logo">S</a>
    <nav class="sidebar-nav">
      <a href="/" class="sidebar-icon" title="Overview"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M2.375 7.125H16.625M2.375 11.875H16.625M7.125 2.375V16.625M11.875 2.375V16.625M3.95833 2.375H15.0417C15.9161 2.375 16.625 3.08388 16.625 3.95833V15.0417C16.625 15.9161 15.9161 16.625 15.0417 16.625H3.95833C3.08388 16.625 2.375 15.9161 2.375 15.0417V3.95833C2.375 3.08388 3.08388 2.375 3.95833 2.375Z"/></svg></a>
      <a href="/secrets" class="sidebar-icon" title="Secrets"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M11.0831 5.54166L13.4582 7.91685M7.44128 8.39158L2.04664 13.7864C1.74967 14.0833 1.58279 14.486 1.5827 14.9059V16.6256C1.5827 16.8356 1.66611 17.0369 1.81458 17.1854C1.96306 17.3339 2.16443 17.4173 2.3744 17.4173H4.7495C4.95947 17.4173 5.16085 17.3339 5.30932 17.1854C5.45779 17.0369 5.5412 16.8356 5.5412 16.6256V15.8338C5.5412 15.6239 5.62461 15.4225 5.77309 15.274C5.92156 15.1255 6.12293 15.0421 6.3329 15.0421H7.1246C7.33458 15.0421 7.53595 14.9587 7.68442 14.8102C7.83289 14.6617 7.91631 14.4604 7.91631 14.2504V13.4587C7.91631 13.2487 7.99972 13.0473 8.14819 12.8988C8.29666 12.7503 8.49803 12.6669 8.70801 12.6669H8.84418C9.26409 12.6668 9.66677 12.4999 9.96364 12.203L10.6081 11.5585M9.81652 2.13736C10.1825 1.79792 10.6632 1.6093 11.1624 1.6093C11.6616 1.6093 12.1423 1.79792 12.5083 2.13736L16.8627 6.49187C17.2021 6.85788 17.3907 7.33863 17.3907 7.83781C17.3907 8.337 17.2021 8.81775 16.8627 9.18376L13.9334 12.1132C13.5674 12.4526 13.0866 12.6412 12.5875 12.6412C12.0883 12.6412 11.6076 12.4526 11.2416 12.1132L6.88723 7.75864C6.5478 7.39263 6.35919 6.91189 6.35919 6.4127C6.35919 5.91352 6.5478 5.43277 6.88723 5.06676L9.81652 2.13736Z"/></svg></a>
      <a href="/providers" class="sidebar-icon" title="Providers"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M16.625 3.95789C16.625 5.26967 13.435 6.33308 9.5 6.33308C5.56497 6.33308 2.375 5.26967 2.375 3.95789M16.625 3.95789C16.625 2.64611 13.435 1.5827 9.5 1.5827C5.56497 1.5827 2.375 2.64611 2.375 3.95789M16.625 3.95789V15.0421C16.625 15.6721 15.8743 16.2762 14.5381 16.7216C13.2019 17.1671 11.3897 17.4173 9.5 17.4173C7.61033 17.4173 5.79806 17.1671 4.46186 16.7216C3.12567 16.2762 2.375 15.6721 2.375 15.0421V3.95789M2.375 9.5C2.375 10.1299 3.12567 10.7341 4.46186 11.1795C5.79806 11.6249 7.61033 11.8752 9.5 11.8752C11.3897 11.8752 13.2019 11.6249 14.5381 11.1795C15.8743 10.7341 16.625 10.1299 16.625 9.5"/></svg></a>
      <a href="/reports" class="sidebar-icon" title="Reports"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M11.0832 1.5827H4.75047C4.33059 1.5827 3.9279 1.74953 3.631 2.04649C3.3341 2.34344 3.1673 2.7462 3.1673 3.16616V15.8338C3.1673 16.2538 3.3341 16.6566 3.631 16.9535C3.9279 17.2505 4.33059 17.4173 4.75047 17.4173H14.2495C14.6694 17.4173 15.0721 17.2505 15.369 16.9535C15.6659 16.6566 15.8327 16.2538 15.8327 15.8338V6.33308M11.0832 1.5827C11.3338 1.5823 11.5819 1.63147 11.8134 1.72741C12.0449 1.82334 12.2552 1.96413 12.432 2.14166L15.2723 4.98239C15.4502 5.15935 15.5914 5.36982 15.6876 5.60165C15.7838 5.83348 15.8331 6.08208 15.8327 6.33308M11.0832 1.5827V5.54135C11.0832 5.75133 11.1666 5.95271 11.315 6.10119C11.4635 6.24966 11.6648 6.33308 11.8748 6.33308L15.8327 6.33308M7.91682 7.12481H6.33365M12.6663 10.2917H6.33365M12.6663 13.4586H6.33365"/></svg></a>
    </nav>
  </div>
  <a href="/settings" class="sidebar-icon active" title="Settings"><svg fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" viewBox="0 0 19 19"><path d="M7.65586 3.2751C7.69948 2.81621 7.91262 2.39006 8.25365 2.07991C8.59467 1.76976 9.03908 1.5979 9.50005 1.5979C9.96101 1.5979 10.4054 1.76976 10.7464 2.07991C11.0875 2.39006 11.3006 2.81621 11.3442 3.2751C11.3705 3.57155 11.4677 3.85731 11.6278 4.10821C11.7878 4.3591 12.006 4.56774 12.2637 4.71647C12.5215 4.8652 12.8113 4.94964 13.1086 4.96264C13.4059 4.97563 13.702 4.91681 13.9718 4.79115C14.3907 4.60097 14.8653 4.57345 15.3034 4.71395C15.7414 4.85444 16.1115 5.1529 16.3416 5.55123C16.5718 5.94956 16.6455 6.41927 16.5484 6.86893C16.4513 7.3186 16.1903 7.71605 15.8164 7.98394C15.5728 8.15482 15.374 8.38184 15.2368 8.64579C15.0995 8.90975 15.0279 9.20287 15.0279 9.50038C15.0279 9.79788 15.0995 10.091 15.2368 10.355C15.374 10.6189 15.5728 10.8459 15.8164 11.0168C16.1903 11.2847 16.4513 11.6822 16.5484 12.1318C16.6455 12.5815 16.5718 13.0512 16.3416 13.4495C16.1115 13.8479 15.7414 14.1463 15.3034 14.2868C14.8653 14.4273 14.3907 14.3998 13.9718 14.2096C13.702 14.0839 13.4059 14.0251 13.1086 14.0381C12.8113 14.0511 12.5215 14.1355 12.2637 14.2843C12.006 14.433 11.7878 14.6416 11.6278 14.8925C11.4677 15.1434 11.3705 15.4292 11.3442 15.7256C11.3006 16.1845 11.0875 16.6107 10.7464 16.9208C10.4054 17.231 9.96101 17.4028 9.50005 17.4028C9.03908 17.4028 8.59467 17.231 8.25365 16.9208C7.91262 16.6107 7.69948 16.1845 7.65586 15.7256C7.62969 15.4291 7.53244 15.1432 7.37234 14.8922C7.21224 14.6413 6.994 14.4325 6.73613 14.2838C6.47825 14.1351 6.18832 14.0507 5.8909 14.0377C5.59348 14.0248 5.29733 14.0838 5.02753 14.2096C4.60865 14.3998 4.13399 14.4273 3.69594 14.2868C3.25789 14.1463 2.88779 13.8479 2.65766 13.4495C2.42753 13.0512 2.35384 12.5815 2.45094 12.1318C2.54803 11.6822 2.80896 11.2847 3.18294 11.0168C3.42648 10.8459 3.62527 10.6189 3.76251 10.355C3.89975 10.091 3.97141 9.79788 3.97141 9.50038C3.97141 9.20287 3.89975 8.90975 3.76251 8.64579C3.62527 8.38184 3.42648 8.15482 3.18294 7.98394C2.80949 7.71591 2.54902 7.31862 2.45216 6.86926C2.3553 6.4199 2.42897 5.95057 2.65885 5.5525C2.88873 5.15443 3.2584 4.85604 3.69601 4.71534C4.13363 4.57463 4.60793 4.60165 5.02673 4.79115C5.2965 4.91681 5.59257 4.97563 5.88989 4.96264C6.18721 4.94964 6.47701 4.8652 6.73478 4.71647C6.99256 4.56774 7.2107 4.3591 7.37075 4.10821C7.53081 3.85731 7.62806 3.57155 7.65428 3.2751"/></svg></a>
</aside>
<main class="main">
  <div class="page-header">
    <h1>Settings</h1>
    <div class="subtitle">Configuration viewer — read only</div>
  </div>
  <div class="content">
    <div class="notice">Edit <code>~/.secretwatch/config.yaml</code> to change settings. Restart <code>secretwatch serve</code> to apply changes.</div>
    <div class="section">
      <h2>Rotation Policy</h2>
      <div class="row"><span class="row-label">Default max age</span><span class="row-val">{{.Cfg.RotationPolicy.DefaultMaxAgeDays}} days</span></div>
      {{range .Cfg.RotationPolicy.Overrides}}
      <div class="row"><span class="row-label">{{.Pattern}}</span><span class="row-val">{{.MaxAgeDays}} days</span></div>
      {{end}}
    </div>
    <div class="section">
      <h2>AWS</h2>
      {{if .Cfg.Stores.AWS}}
      {{range .Cfg.Stores.AWS}}
      <div class="row"><span class="row-label">Profile</span><span class="row-val">{{if .Profile}}{{.Profile}}{{else}}default{{end}}</span></div>
      <div class="row"><span class="row-label">Regions</span><span class="row-val">{{range .Regions}}{{.}} {{end}}</span></div>
      {{end}}
      {{else}}<div class="row"><span class="row-label">Status</span><span class="row-val muted">not configured</span></div>{{end}}
    </div>
    <div class="section">
      <h2>GCP</h2>
      {{if .Cfg.Stores.GCP}}
      {{range .Cfg.Stores.GCP}}
      <div class="row"><span class="row-label">Project</span><span class="row-val">{{.Project}}</span></div>
      {{end}}
      {{else}}<div class="row"><span class="row-label">Status</span><span class="row-val muted">not configured</span></div>{{end}}
    </div>
    <div class="section">
      <h2>Vault</h2>
      <div class="row"><span class="row-label">Status</span><span>{{if .Cfg.Stores.Vault.Enabled}}<span class="badge-on">enabled</span>{{else}}<span class="badge-off">disabled</span>{{end}}</span></div>
      {{if .Cfg.Stores.Vault.Enabled}}
      <div class="row"><span class="row-label">Address</span><span class="row-val">{{.Cfg.Stores.Vault.Address}}</span></div>
      {{end}}
    </div>
    <div class="section">
      <h2>GitHub</h2>
      <div class="row"><span class="row-label">Status</span><span>{{if .Cfg.Stores.GitHub.Enabled}}<span class="badge-on">enabled</span>{{else}}<span class="badge-off">disabled</span>{{end}}</span></div>
      {{if .Cfg.Stores.GitHub.Enabled}}
      <div class="row"><span class="row-label">Org</span><span class="row-val">{{.Cfg.Stores.GitHub.Org}}</span></div>
      {{end}}
    </div>
    <div class="section">
      <h2>Storage</h2>
      <div class="row"><span class="row-label">DB path</span><span class="row-val">{{.Cfg.DBPath}}</span></div>
      <div class="row"><span class="row-label">Web port</span><span class="row-val">{{.Cfg.Web.Port}}</span></div>
    </div>
  </div>
</main>
</body></html>
`
