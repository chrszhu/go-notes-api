package notes

import (
	"context"
	"database/sql"
	"errors"
)

type Note struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Repository interface {
	Create(ctx context.Context, n *Note) error
	Get(ctx context.Context, id int64) (*Note, error)
	List(ctx context.Context) ([]Note, error)
}

type repo struct { db *sql.DB }

func NewRepository(db *sql.DB) Repository {
	return &repo{db: db}
}

func (r *repo) ensureTable() error {
	_, err := r.db.Exec(`CREATE TABLE IF NOT EXISTS notes (id SERIAL PRIMARY KEY, title TEXT NOT NULL, content TEXT NOT NULL)`)
	return err
}

func (r *repo) Create(ctx context.Context, n *Note) error {
	if err := r.ensureTable(); err != nil { return err }
	return r.db.QueryRowContext(ctx, `INSERT INTO notes (title, content) VALUES ($1, $2) RETURNING id`, n.Title, n.Content).Scan(&n.ID)
}

func (r *repo) Get(ctx context.Context, id int64) (*Note, error) {
	if err := r.ensureTable(); err != nil { return nil, err }
	var n Note
	row := r.db.QueryRowContext(ctx, `SELECT id, title, content FROM notes WHERE id=$1`, id)
	if err := row.Scan(&n.ID, &n.Title, &n.Content); err != nil {
		if errors.Is(err, sql.ErrNoRows) { return nil, nil }
		return nil, err
	}
	return &n, nil
}

func (r *repo) List(ctx context.Context) ([]Note, error) {
	if err := r.ensureTable(); err != nil { return nil, err }
	rows, err := r.db.QueryContext(ctx, `SELECT id, title, content FROM notes ORDER BY id DESC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var result []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content); err != nil { return nil, err }
		result = append(result, n)
	}
	return result, rows.Err()
}
