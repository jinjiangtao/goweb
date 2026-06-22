package store

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"qrcode/internal/model"
)

type Store struct {
	DB *sqlx.DB
}

var schema = `
CREATE TABLE IF NOT EXISTS templates (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    name         TEXT    NOT NULL,
    style_config TEXT    NOT NULL,
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS records (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    content      TEXT    NOT NULL,
    content_type TEXT    NOT NULL DEFAULT 'text',
    source       TEXT    NOT NULL DEFAULT 'single',
    style_config TEXT    NOT NULL,
    file_path    TEXT    NOT NULL,
    status       TEXT    NOT NULL DEFAULT 'active',
    template_id  INTEGER,
    batch_id     TEXT,
    remark       TEXT,
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (template_id) REFERENCES templates(id)
);

CREATE INDEX IF NOT EXISTS idx_records_status  ON records(status);
CREATE INDEX IF NOT EXISTS idx_records_batch   ON records(batch_id);
CREATE INDEX IF NOT EXISTS idx_records_created ON records(created_at);

CREATE TABLE IF NOT EXISTS parse_logs (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    source     TEXT,
    result     TEXT,
    type       TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

func Open(dbPath string) (*Store, error) {
	db, err := sqlx.Connect("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return &Store{DB: db}, nil
}

func (s *Store) CreateRecord(r *model.Record) error {
	res, err := s.DB.NamedExec(
		`INSERT INTO records (content, content_type, source, style_config, file_path, status, template_id, batch_id, remark)
		 VALUES (:content, :content_type, :source, :style_config, :file_path, :status, :template_id, :batch_id, :remark)`,
		map[string]interface{}{
			"content":      r.Content,
			"content_type": r.ContentType,
			"source":       r.Source,
			"style_config": r.StyleConfig.JSON(),
			"file_path":    r.FilePath,
			"status":       r.Status,
			"template_id":  r.TemplateID,
			"batch_id":     r.BatchID,
			"remark":       r.Remark,
		},
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	r.ID = id
	return s.DB.Get(r, `SELECT * FROM records WHERE id = ?`, id)
}

func (s *Store) GetRecord(id int64) (model.Record, error) {
	var r model.Record
	err := s.DB.Get(&r, `SELECT * FROM records WHERE id = ?`, id)
	return r, err
}

func (s *Store) ListRecords(status, keyword string, page, size int) ([]model.Record, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 200 {
		size = 20
	}
	where := []string{"1=1"}
	args := []interface{}{}
	if status != "" {
		where = append(where, "status = ?")
		args = append(args, status)
	}
	if keyword != "" {
		where = append(where, "(content LIKE ? OR batch_id LIKE ?)")
		args = append(args, "%"+keyword+"%", "%"+keyword+"%")
	}
	wsql := strings.Join(where, " AND ")

	var total int64
	if err := s.DB.Get(&total, `SELECT COUNT(*) FROM records WHERE `+wsql, args...); err != nil {
		return nil, 0, err
	}

	args = append(args, size, (page-1)*size)
	var list []model.Record
	err := s.DB.Select(&list,
		`SELECT * FROM records WHERE `+wsql+` ORDER BY id DESC LIMIT ? OFFSET ?`, args...)
	return list, total, err
}

func (s *Store) UpdateRecordStatus(id int64, status string) error {
	_, err := s.DB.Exec(
		`UPDATE records SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, id)
	return err
}

func (s *Store) UpdateRecordContent(id int64, content, contentType, filePath string) error {
	_, err := s.DB.Exec(
		`UPDATE records SET content = ?, content_type = ?, file_path = ?, status = 'active', updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		content, contentType, filePath, id)
	return err
}

func (s *Store) DeleteRecord(id int64) error {
	_, err := s.DB.Exec(`DELETE FROM records WHERE id = ?`, id)
	return err
}

func (s *Store) GetRecordsByIDs(ids []int64) ([]model.Record, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	q := `SELECT * FROM records WHERE id IN (` + strings.Join(placeholders, ",") + `)`
	var list []model.Record
	err := s.DB.Select(&list, q, args...)
	return list, err
}

func (s *Store) CreateTemplate(t *model.Template) error {
	res, err := s.DB.NamedExec(
		`INSERT INTO templates (name, style_config) VALUES (:name, :style_config)`,
		map[string]interface{}{"name": t.Name, "style_config": t.StyleConfig.JSON()},
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	t.ID = id
	return s.DB.Get(t, `SELECT * FROM templates WHERE id = ?`, id)
}

func (s *Store) GetTemplate(id int64) (model.Template, error) {
	var t model.Template
	err := s.DB.Get(&t, `SELECT * FROM templates WHERE id = ?`, id)
	return t, err
}

func (s *Store) ListTemplates(keyword string) ([]model.Template, error) {
	var list []model.Template
	if keyword == "" {
		err := s.DB.Select(&list, `SELECT * FROM templates ORDER BY id DESC`)
		return list, err
	}
	err := s.DB.Select(&list, `SELECT * FROM templates WHERE name LIKE ? ORDER BY id DESC`, "%"+keyword+"%")
	return list, err
}

func (s *Store) UpdateTemplate(id int64, name string, style model.StyleConfig) error {
	_, err := s.DB.Exec(
		`UPDATE templates SET name = ?, style_config = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		name, style.JSON(), id)
	return err
}

func (s *Store) DeleteTemplate(id int64) error {
	_, err := s.DB.Exec(`DELETE FROM templates WHERE id = ?`, id)
	return err
}

func (s *Store) CreateParseLog(source, result, typ string) {
	_, _ = s.DB.Exec(
		`INSERT INTO parse_logs (source, result, type) VALUES (?, ?, ?)`, source, result, typ)
}

func (s *Store) Stats() (model.Stats, error) {
	var st model.Stats
	if err := s.DB.Get(&st.Total, `SELECT COUNT(*) FROM records`); err != nil {
		return st, err
	}
	if err := s.DB.Get(&st.Today, `SELECT COUNT(*) FROM records WHERE date(created_at) = date('now')`); err != nil {
		return st, err
	}
	if err := s.DB.Get(&st.Templates, `SELECT COUNT(*) FROM templates`); err != nil {
		return st, err
	}
	if err := s.DB.Get(&st.Invalid, `SELECT COUNT(*) FROM records WHERE status = 'invalid'`); err != nil {
		return st, err
	}
	var recent []model.Record
	if err := s.DB.Select(&recent, `SELECT * FROM records ORDER BY id DESC LIMIT 8`); err != nil {
		return st, err
	}
	st.Recent = recent
	if st.Recent == nil {
		st.Recent = []model.Record{}
	}
	return st, nil
}
