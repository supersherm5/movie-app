package postgres

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/supersherm5/movie-app/rating/internal/repository"
	"github.com/supersherm5/movie-app/rating/pkg/model"
)

// Repository defines a Postgres-based movie metadata repo
type Repository struct {
	db *sql.DB
}

// New creates a new Postgres repository
func NewPostgresRepo() (*Repository, error) {
	dsn := "user=postgres password=mysecretpassword dbname=movies host=localhost port=5432 dbname=movies sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

// Get retrieves all ratings for a given record.
func (r *Repository) Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT user_id, value FROM ratings WHERE record_id = $1 AND record_type = $2", recordID, recordType)
	if err != nil {
			return nil, err
	}

	defer rows.Close()
	var res []model.Rating
	for rows.Next() {
			var userID string
			var value int32
			if err := rows.Scan(&userID, &value); err != nil {
					return nil, err
			}
			res = append(res, model.Rating{
					UserID: model.UserID(userID),
					Value:  model.RatingValue(value),
			})
	}
	if len(res) == 0 {
			return nil, repository.ErrNoRatingExists
	}
	return res, nil
}
// Put adds a rating for a given record.
func (r *Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO ratings (record_id, record_type, user_id, value) VALUES ($1, $2, $3, $4)", recordID, recordType, rating.UserID, rating.Value)
		if err != nil {
			log.Printf("[Rating] Put Err: %v\n", err)
		}
	return err
}