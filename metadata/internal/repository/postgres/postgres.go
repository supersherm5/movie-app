package postgres

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/supersherm5/movie-app/metadata/internal/repository"
	"github.com/supersherm5/movie-app/metadata/pkg/model"
)

// Repository defines a Postgres-based movie metadata repo
type Repository struct {
	db *sql.DB
}

// New creates a new Postgres repository
func NewPostgresRepo() (*Repository, error) {
	connStr := "user=postgres password=mysecretpassword dbname=movies host=localhost port=5432 dbname=movies sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

// Get retrieves moviee metadata for by movie id.
func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	var title, description, director string
	row := r.db.QueryRowContext(ctx, "Select title, description, director FROM metadata WHERE id = $1", id)
	if err := row.Scan(&title, &description, &director); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrMetadataNotFound
		}
		return nil, err
	}

	return &model.Metadata{
		ID: id,
		Title: title,
		Description: description,
		Director: director,
	}, nil
}

// Put adds movie metadata for a given movie id.
func (r *Repository) Put(ctx context.Context, id string, metadata *model.Metadata) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO metadata (id, title, description, director) VALUES ($1, $2, $3, $4)", id, metadata.Title, metadata.Description, metadata.Director)
	
	return err
}