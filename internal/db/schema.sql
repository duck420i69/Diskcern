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
