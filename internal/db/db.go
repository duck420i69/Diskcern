package db

import (
	"database/sql"
	"fmt"

	"github.com/diskcern/diskcern/internal/models"
	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

func InitDB(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Load schema
	// For simplicity in this step, we read a fixed schema file or embed it.
	// Assume schema is already executed or we can embed it later.
	// But let's create the tables directly here for robustness.
	schema := `
	CREATE TABLE IF NOT EXISTS snapshots (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		root_path TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS file_records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		snapshot_id INTEGER NOT NULL,
		path TEXT NOT NULL,
		size INTEGER NOT NULL,
		is_dir BOOLEAN NOT NULL CHECK (is_dir IN (0, 1)),
		matched_rule TEXT,
		description TEXT,
		FOREIGN KEY(snapshot_id) REFERENCES snapshots(id) ON DELETE CASCADE
	);
	
	CREATE INDEX IF NOT EXISTS idx_file_records_snapshot_path ON file_records(snapshot_id, path);

	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);

	CREATE TABLE IF NOT EXISTS file_record_tags (
		file_record_id INTEGER NOT NULL,
		tag_id INTEGER NOT NULL,
		PRIMARY KEY(file_record_id, tag_id),
		FOREIGN KEY(file_record_id) REFERENCES file_records(id) ON DELETE CASCADE,
		FOREIGN KEY(tag_id) REFERENCES tags(id) ON DELETE CASCADE
	);
	`
	_, err = conn.Exec(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return &DB{conn: conn}, nil
}

func (d *DB) CreateSnapshot(rootPath string) (int64, error) {
	res, err := d.conn.Exec("INSERT INTO snapshots (root_path) VALUES (?)", rootPath)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (d *DB) InsertRecords(snapshotID int64, records []models.FileRecord) error {
	tx, err := d.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO file_records (snapshot_id, path, size, is_dir, matched_rule, description) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	tagStmt, err := tx.Prepare("INSERT OR IGNORE INTO tags (name) VALUES (?)")
	if err != nil {
		return err
	}
	defer tagStmt.Close()

	tagIDStmt, err := tx.Prepare("SELECT id FROM tags WHERE name = ?")
	if err != nil {
		return err
	}
	defer tagIDStmt.Close()

	linkStmt, err := tx.Prepare("INSERT INTO file_record_tags (file_record_id, tag_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer linkStmt.Close()

	for _, r := range records {
		res, err := stmt.Exec(snapshotID, r.Path, r.Size, r.IsDir, r.MatchedRule, r.Description)
		if err != nil {
			return err
		}
		recordID, err := res.LastInsertId()
		if err != nil {
			return err
		}

		for _, tag := range r.Tags {
			_, err = tagStmt.Exec(tag)
			if err != nil {
				return err
			}
			var tagID int64
			err = tagIDStmt.QueryRow(tag).Scan(&tagID)
			if err != nil {
				return err
			}
			_, err = linkStmt.Exec(recordID, tagID)
			if err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

func (d *DB) Close() error {
	return d.conn.Close()
}
