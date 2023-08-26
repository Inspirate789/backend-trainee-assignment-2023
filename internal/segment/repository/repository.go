package repository

import (
	"context"
	"database/sql"
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/segment/usecase"
	"github.com/Inspirate789/backend-trainee-assignment-2023/pkg/sqlx_utils"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"time"
)

type sqlxSegmentRepository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewSqlxSegmentRepository(db *sqlx.DB, logger *slog.Logger) usecase.SegmentRepository {
	return &sqlxSegmentRepository{
		db:     db,
		logger: logger.WithGroup("sqlxSegmentRepository"),
	}
}

func (r *sqlxSegmentRepository) AddSegment(name string, userPercentage float64, ttl time.Duration) error {
	args := map[string]any{
		"seg_name":        name,
		"user_percentage": userPercentage,
		"expire": sql.NullTime{
			Time:  time.Now().Add(ttl),
			Valid: ttl != usecase.NoTTL,
		},
	}

	_, err := sqlx_utils.NamedExec(context.Background(), r.db, insertSegmentQuery, args)

	return err
}

func (r *sqlxSegmentRepository) RemoveSegment(name string) error {
	args := map[string]any{
		"seg_name": name,
	}
	_, err := sqlx_utils.NamedExec(context.Background(), r.db, deleteSegmentQuery, args)

	return err
}
