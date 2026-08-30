package ledger

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	_ "modernc.org/sqlite"
)

type Record struct {
	Workspace string
	Cloid     string
	Preview   string
	Status    string
	OID       string
}

type Store struct {
	db        *sql.DB
	mu        sync.Mutex
	workspace string
}

func Open(dir, network, workspaceID string) (*Store, error) {
	if workspaceID == "" {
		return nil, fmt.Errorf("workspace required")
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	path := filepath.Join(dir, network+"-"+workspaceID+".db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS actions (
		cloid TEXT PRIMARY KEY,
		workspace TEXT NOT NULL,
		preview TEXT NOT NULL,
		status TEXT NOT NULL,
		oid TEXT
	)`); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Store{db: db, workspace: workspaceID}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) Apply(rec Record) (applied bool, err error) {
	if err := s.RefuseForeign(rec.Workspace); err != nil {
		return false, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	res, err := s.db.Exec(`INSERT INTO actions(cloid, workspace, preview, status, oid) VALUES(?,?,?,?,?)`,
		rec.Cloid, rec.Workspace, rec.Preview, rec.Status, rec.OID)
	if err != nil {
		if isUnique(err) {
			return false, nil
		}
		return false, err
	}
	n, _ := res.RowsAffected()
	return n == 1, nil
}

func (s *Store) Get(workspace, cloid string) (Record, error) {
	if err := s.RefuseForeign(workspace); err != nil {
		return Record{}, err
	}
	var r Record
	err := s.db.QueryRow(`SELECT cloid, workspace, preview, status, oid FROM actions WHERE cloid=? AND workspace=?`, cloid, workspace).
		Scan(&r.Cloid, &r.Workspace, &r.Preview, &r.Status, &r.OID)
	if err == sql.ErrNoRows {
		return Record{}, fmt.Errorf("not found")
	}
	return r, err
}

func (s *Store) HasPreview(workspace, preview string) bool {
	if strings.TrimSpace(preview) == "" {
		return false
	}
	if err := s.RefuseForeign(workspace); err != nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	var n int
	err := s.db.QueryRow(`SELECT COUNT(1) FROM actions WHERE workspace=? AND preview=?`, workspace, preview).Scan(&n)
	return err == nil && n > 0
}

func isUnique(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "UNIQUE") || strings.Contains(msg, "unique")
}
