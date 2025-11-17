package domain_service

import (
	"context"
	"log/slog"

	domain "github.com/example/avito_test_task_internship/internal/domain/entity"
	"github.com/example/avito_test_task_internship/internal/infrasctucture"
	"github.com/jackc/pgx/v5"
)

type UserService struct {
	database *infrasctucture.PostgresDb
}

func (s *UserService) SetIsActive(ctx context.Context, userId string, isActive bool) (*domain.User, error) {
	tx, err := s.database.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		slog.Error("transacton Set Is Active failed", "error", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	var u domain.User
	err = tx.QueryRow(ctx,
		`SELECT user_id, username, is_active, team_name 
		 FROM users 
		 WHERE user_id = $1 
		 FOR UPDATE`, userId).Scan(&u.UserID, &u.Username, &u.IsActive, &u.TeamName)
	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Error("user not found", "user_id", userId)
			return nil, ErrUserNotFound
		}
		slog.Error("failed to select user", "error", err)
		return nil, err
	}

	_, err = tx.Exec(ctx,
		`UPDATE users SET is_active = $1 WHERE user_id = $2`,
		isActive, userId)
	if err != nil {
		slog.Error("transacton Set Is Active failed", "error", err)
		return nil, err
	}

	u.IsActive = isActive

	if err := tx.Commit(ctx); err != nil {
		slog.Error("transacton Set Is Active failed", "error", err)
		return nil, err
	}

	return &u, nil
}

func (s *UserService) GetReviews(ctx context.Context, userId string) ([]*domain.PullRequest, error) {
	rows, err := s.database.Pool.Query(ctx,
		`SELECT id, name, status, author, createdat, mergedat 
		 FROM pullRequests 
		 WHERE author = $1`, userId)
	if err != nil {
		slog.Error("query pull Get Review failed", "error", err)
		return nil, err
	}
	defer rows.Close()

	prs := make([]*domain.PullRequest, 0)
	for rows.Next() {
		var pr domain.PullRequest
		if err := rows.Scan(&pr.PullRequestId, &pr.Name, &pr.Status, &pr.AuthorId, &pr.CreatedAt, &pr.MergedAt); err != nil {
			slog.Error("query pull Get Review failed", "error", err)
			return nil, err
		}
		prs = append(prs, &pr)
	}

	if len(prs) == 0 {
		slog.Info("no pull requests found for user:", "userId", userId)
		return nil, ErrNoPullRequestUser
	}

	return prs, nil
}

func NewUserService(database *infrasctucture.PostgresDb) *UserService {
	return &UserService{database}
}
