package repository

import (
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/segments/usecase"
	"github.com/jmoiron/sqlx"
	"log/slog"
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

func (s *sqlxSegmentRepository) AddSegment(name string) error {
	//TODO implement me
	panic("implement me")
}

func (s *sqlxSegmentRepository) RemoveSegment(name string) error {
	//TODO implement me
	panic("implement me")
}
