package db

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
	"github.com/leylmordor/secretwatch/internal/store"
)

type DB struct {
	db *sql.DB
}

func Open(path string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }

func migrate(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS secrets (
			full_path    TEXT PRIMARY KEY,
			name         TEXT NOT NULL,
			store_type   TEXT NOT NULL,
			region       TEXT,
			last_rotated DATETIME,
			max_age_days INTEGER NOT NULL,
			status       TEXT NOT NULL,
			days_until_expiry INTEGER,
			tags         TEXT,
			scanned_at   DATETIME NOT NULL
		)
	`); err != nil {
		return err
	}
	_, _ = db.Exec(`ALTER TABLE secrets ADD COLUMN excluded INTEGER NOT NULL DEFAULT 0`)
	_, _ = db.Exec(`ALTER TABLE secrets ADD COLUMN account TEXT NOT NULL DEFAULT ''`)
	return nil
}

func (d *DB) Save(secrets []*store.Secret) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().UTC().Format(time.RFC3339)
	for _, s := range secrets {
		tags, _ := json.Marshal(s.Tags)

		var lastRotated *time.Time
		if s.LastRotated != nil {
			lastRotated = s.LastRotated
		}

		var daysUntil *int
		if s.DaysUntilExpiry != nil {
			daysUntil = s.DaysUntilExpiry
		}

		_, err := tx.Exec(`
			INSERT INTO secrets (full_path, name, store_type, region, last_rotated, max_age_days, status, days_until_expiry, tags, scanned_at, account)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(full_path) DO UPDATE SET
				name=excluded.name,
				store_type=excluded.store_type,
				region=excluded.region,
				last_rotated=excluded.last_rotated,
				max_age_days=excluded.max_age_days,
				status=excluded.status,
				days_until_expiry=excluded.days_until_expiry,
				tags=excluded.tags,
				scanned_at=excluded.scanned_at,
				account=excluded.account
		`, s.FullPath, s.Name, s.StoreType, s.Region, lastRotated, s.MaxAgeDays, string(s.Status), daysUntil, string(tags), now, s.Account)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func scanSecret(rows *sql.Rows) (*store.Secret, error) {
	s := &store.Secret{Tags: make(map[string]string)}
	var lastRotated sql.NullTime
	var region sql.NullString
	var daysUntil sql.NullInt64
	var tagsJSON string
	var excluded int
	var account sql.NullString

	if err := rows.Scan(&s.FullPath, &s.Name, &s.StoreType, &region, &lastRotated, &s.MaxAgeDays, &s.Status, &daysUntil, &tagsJSON, &excluded, &account); err != nil {
		return nil, err
	}
	if lastRotated.Valid {
		s.LastRotated = &lastRotated.Time
	}
	if region.Valid {
		s.Region = region.String
	}
	if daysUntil.Valid {
		v := int(daysUntil.Int64)
		s.DaysUntilExpiry = &v
	}
	if account.Valid {
		s.Account = account.String
	}
	_ = json.Unmarshal([]byte(tagsJSON), &s.Tags)
	s.Excluded = excluded == 1
	return s, nil
}

type StoreAccountPair struct {
	StoreType string
	Account   string
	Count     int
}

func (d *DB) StoreAccountPairs() ([]StoreAccountPair, error) {
	rows, err := d.db.Query(`
		SELECT store_type, COALESCE(account,'') as account, COUNT(*) as cnt
		FROM secrets WHERE excluded=0
		GROUP BY store_type, account ORDER BY store_type, account
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pairs []StoreAccountPair
	for rows.Next() {
		var p StoreAccountPair
		if err := rows.Scan(&p.StoreType, &p.Account, &p.Count); err != nil {
			return nil, err
		}
		pairs = append(pairs, p)
	}
	return pairs, rows.Err()
}

func (d *DB) List(storeType, accountFilter, statusFilter string, showExcluded bool) ([]*store.Secret, error) {
	query := `SELECT full_path, name, store_type, region, last_rotated, max_age_days, status, days_until_expiry, tags, excluded, account FROM secrets WHERE 1=1`
	var args []interface{}

	if !showExcluded {
		query += " AND excluded = 0"
	}
	if storeType != "" {
		query += " AND store_type = ?"
		args = append(args, storeType)
	}
	if accountFilter != "" {
		query += " AND account = ?"
		args = append(args, accountFilter)
	}
	if statusFilter != "" {
		query += " AND status = ?"
		args = append(args, statusFilter)
	}
	query += " ORDER BY excluded ASC, status DESC, name ASC"

	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var secrets []*store.Secret
	for rows.Next() {
		s, err := scanSecret(rows)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, s)
	}
	return secrets, rows.Err()
}

func (d *DB) RecentRotations(limit int) ([]*store.Secret, error) {
	rows, err := d.db.Query(`SELECT full_path, name, store_type, region, last_rotated, max_age_days, status, days_until_expiry, tags, excluded, account FROM secrets WHERE last_rotated IS NOT NULL AND excluded = 0 ORDER BY last_rotated DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var secrets []*store.Secret
	for rows.Next() {
		s, err := scanSecret(rows)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, s)
	}
	return secrets, rows.Err()
}

func (d *DB) SetExcluded(fullPath string, excluded bool) error {
	v := 0
	if excluded {
		v = 1
	}
	_, err := d.db.Exec(`UPDATE secrets SET excluded = ? WHERE full_path = ?`, v, fullPath)
	return err
}

func (d *DB) StoreTypes() []string {
	rows, err := d.db.Query(`SELECT DISTINCT store_type FROM secrets ORDER BY store_type`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var types []string
	for rows.Next() {
		var t string
		if rows.Scan(&t) == nil {
			types = append(types, t)
		}
	}
	return types
}

func (d *DB) GetStoreType(name string) string {
	var storeType string
	d.db.QueryRow(`SELECT store_type FROM secrets WHERE name = ? LIMIT 1`, name).Scan(&storeType)
	return storeType
}

type StoreStat struct {
	StoreType string
	Total     int
	Overdue   int
	Warning   int
	OK        int
	Unknown   int
}

func (d *DB) StoreStats() ([]StoreStat, error) {
	rows, err := d.db.Query(`
		SELECT store_type,
		       COUNT(*) as total,
		       SUM(CASE WHEN status='overdue' THEN 1 ELSE 0 END) as overdue,
		       SUM(CASE WHEN status='warning' THEN 1 ELSE 0 END) as warning,
		       SUM(CASE WHEN status='ok' THEN 1 ELSE 0 END) as ok,
		       SUM(CASE WHEN status='unknown' THEN 1 ELSE 0 END) as unknown
		FROM secrets WHERE excluded=0
		GROUP BY store_type ORDER BY store_type
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var stats []StoreStat
	for rows.Next() {
		var s StoreStat
		if err := rows.Scan(&s.StoreType, &s.Total, &s.Overdue, &s.Warning, &s.OK, &s.Unknown); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, rows.Err()
}

func (d *DB) LastScan() (*time.Time, error) {
	var t sql.NullTime
	err := d.db.QueryRow(`SELECT MAX(scanned_at) FROM secrets`).Scan(&t)
	if err != nil || !t.Valid {
		return nil, err
	}
	return &t.Time, nil
}
