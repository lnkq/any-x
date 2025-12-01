package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"any-x/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s := &Storage{db: db}
	if err := s.init(); err != nil {
		return nil, fmt.Errorf("%s: failed to init db: %w", op, err)
	}

	return s, nil
}

func (s *Storage) init() error {
	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		created_at DATETIME NOT NULL
	);
	`
	_, err := s.db.Exec(query)
	return err
}

func (s *Storage) SavePost(ctx context.Context, post *models.Post) error {
	const op = "storage.sqlite.SavePost"
	query := "INSERT INTO posts (content, created_at) VALUES (?, ?)"
	res, err := s.db.ExecContext(ctx, query, post.Content, post.CreatedAt)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	post.ID = id
	return nil
}

func (s *Storage) GetPosts(ctx context.Context, limit int) ([]models.Post, error) {
	const op = "storage.sqlite.GetPosts"
	query := "SELECT id, content, created_at FROM posts ORDER BY created_at DESC LIMIT ?"
	rows, err := s.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.Content, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return posts, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
