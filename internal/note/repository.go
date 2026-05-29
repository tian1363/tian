package note

import (
	"context"
	"database/sql"
	"errors"
)

var ErrNoteNotFound = errors.New("note not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, req CreateNoteRequest) (int64, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO notes (title, content) VALUES (?, ?)`, req.Title, req.Content)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) List(ctx context.Context) ([]Note, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, content, created_at, updated_at
		FROM notes
		ORDER BY updated_at DESC, id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]Note, 0)
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (Note, error) {
	var n Note
	err := r.db.QueryRowContext(ctx, `
		SELECT id, title, content, created_at, updated_at
		FROM notes
		WHERE id = ?`, id).
		Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Note{}, ErrNoteNotFound
	}
	if err != nil {
		return Note{}, err
	}
	return n, nil
}

func (r *Repository) Update(ctx context.Context, id int64, req UpdateNoteRequest) error {
	res, err := r.db.ExecContext(ctx, `UPDATE notes SET title = ?, content = ? WHERE id = ?`, req.Title, req.Content, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNoteNotFound
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM notes WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNoteNotFound
	}
	return nil
}
