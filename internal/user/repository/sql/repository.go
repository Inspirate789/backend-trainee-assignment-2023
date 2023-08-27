package sql

import (
	"context"
	"database/sql"
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/models"
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/segment/usecase/dto"
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/user/usecase"
	"github.com/Inspirate789/backend-trainee-assignment-2023/pkg/sqlx_utils"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"log/slog"
	"time"
)

type sqlxPgUserRepository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewSqlxRepository(db *sqlx.DB, logger *slog.Logger) usecase.SqlRepository {
	return &sqlxPgUserRepository{
		db:     db,
		logger: logger.WithGroup("sqlxPgUserRepository"),
	}
}

func (r *sqlxPgUserRepository) AddUser(userID int) error {
	args := map[string]any{
		"user_id": userID,
	}
	_, err := sqlx_utils.NamedExec(context.Background(), r.db, insertUserQuery, args)

	return err
}

func (r *sqlxPgUserRepository) RemoveUser(userID int) error {
	args := map[string]any{
		"user_id": userID,
	}
	_, err := sqlx_utils.NamedExec(context.Background(), r.db, deleteUserQuery, args)

	return err
}

func (r *sqlxPgUserRepository) addUserSegments(tx *sqlx.Tx, userID int, segmentNames []string, ttl time.Duration) error {
	args := map[string]any{
		"user_id": userID,
		"names":   pq.Array(segmentNames),
		"expire": sql.NullTime{
			Time:  time.Now().Add(ttl),
			Valid: ttl != dto.NoTTL,
		},
	}
	_, err := sqlx_utils.NamedExec(context.Background(), tx, insertUserSegmentsQuery, args)

	return err
}

func (r *sqlxPgUserRepository) removeUserSegments(tx *sqlx.Tx, userID int, segmentNames []string) error {
	args := map[string]any{
		"user_id": userID,
		"names":   pq.Array(segmentNames),
	}
	_, err := sqlx_utils.NamedExec(context.Background(), tx, deleteUserSegmentsQuery, args)

	return err
}

func (r *sqlxPgUserRepository) changeUserSegmentsTx(tx *sqlx.Tx, userID int, oldSegmentNames, newSegmentNames []string, ttl time.Duration) error {
	r.logger.Debug("begin transaction")
	r.logger.Debug("begin transaction",
		slog.Int("user_id", userID),
		"old_segment_names", oldSegmentNames,
		"new_segment_names", newSegmentNames,
	)
	err := r.removeUserSegments(tx, userID, oldSegmentNames)
	if err != nil {
		return errors.Wrap(err, "cannot remove old user segments")
	}

	err = r.addUserSegments(tx, userID, newSegmentNames, ttl)
	if err != nil {
		return errors.Wrap(err, "cannot add new user segments")
	}
	r.logger.Debug("transaction completed successfully")

	return nil
}

func (r *sqlxPgUserRepository) ChangeUserSegments(userID int, oldSegmentNames, newSegmentNames []string, ttl time.Duration) error {
	return sqlx_utils.RunTx(context.Background(), r.db, func(tx *sqlx.Tx) error {
		err := r.changeUserSegmentsTx(tx, userID, oldSegmentNames, newSegmentNames, ttl)
		return err
	})
}

func (r *sqlxPgUserRepository) GetUserSegments(userID int) ([]string, error) {
	args := map[string]any{
		"user_id": userID,
	}
	segments := make([]string, 0)
	err := sqlx_utils.NamedSelect(context.Background(), r.db, &segments, selectUserSegmentsQuery, args)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return segments, nil
}

func (r *sqlxPgUserRepository) GetUserHistory(year, month int) ([]models.SegmentEvent, error) {
	args := map[string]any{
		"year":  year,
		"month": month,
	}
	events := make([]models.SegmentEvent, 0)
	err := sqlx_utils.NamedSelect(context.Background(), r.db, &events, selectUserHistoryQuery, args)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return events, nil
}
