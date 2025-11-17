package domain_service

import (
	"context"
	"log/slog"
	"time"

	domain "github.com/example/avito_test_task_internship/internal/domain/entity"
	"github.com/example/avito_test_task_internship/internal/infrasctucture"
	"github.com/jackc/pgx/v5"
)

type PullRequestService struct {
	database *infrasctucture.PostgresDb
}

func (s *PullRequestService) Create(ctx context.Context, pullRequestId string, authorId string, name string) (*domain.PullRequest, error) {
	tx, err := s.database.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		slog.Error("transacton create pull request failed", "error", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	var author_exists bool
	err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)`, authorId).Scan(&author_exists)
	if err != nil {
		slog.Error("failed to check author existence", "error", err)
		return nil, err
	}
	if !author_exists {
		slog.Info("author not found", "authorId", authorId)
		return nil, ErrUserNotFound
	}

	var count int
	err = tx.QueryRow(ctx, `SELECT COUNT(*) FROM pullRequests WHERE id = $1`, pullRequestId).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		slog.Info("pull request already exists with ID", "pullRequestId", pullRequestId)
		return nil, ErrPullRequestAlreadyExists
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO pullRequests (id, name, status, author, CreatedAt) VALUES ($1, $2, 'OPEN', $3, NOW())`,
		pullRequestId, name, authorId)
	if err != nil {
		slog.Error("transacton create pull request failed", "error", err)
		return nil, err
	}

	pr := &domain.PullRequest{
		PullRequestId: pullRequestId,
		Name:          name,
		Status:        "OPEN",
		AuthorId:      authorId,
		CreatedAt:     time.Now(),
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error("transacton create pull request failed", "error", err)
		return nil, err
	}

	return pr, nil
}

// Merge меняет статус PullRequest на MERGED и автоматически выставляет MergedAt
func (s *PullRequestService) Merge(ctx context.Context, pullRequestId string) (*domain.PullRequest, error) {
	tx, err := s.database.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		slog.Error("transacton merge pull request failed", "error", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	var pr domain.PullRequest
	err = tx.QueryRow(ctx,
		`SELECT id, name, status, author, CreatedAt, MergedAt FROM pullRequests WHERE id = $1 FOR UPDATE`,
		pullRequestId).Scan(&pr.PullRequestId, &pr.Name, &pr.Status, &pr.AuthorId, &pr.CreatedAt, &pr.MergedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Info("pull request not found:", pullRequestId, pullRequestId)
			return nil, ErrPullRequestNotFound
		}
		return nil, err
	}

	_, err = tx.Exec(ctx, `UPDATE pullRequests SET status = 'MERGED', MergedAt = NOW() WHERE id = $1`, pullRequestId)
	if err != nil {
		return nil, err
	}

	pr.Status = "MERGED"

	if err := tx.Commit(ctx); err != nil {
		slog.Error("transacton merge pull request failed", "error", err)
		return nil, err
	}

	return &pr, nil
}

func (s *PullRequestService) Reassign(ctx context.Context, pullRequestId string, oldReviewerId string) (*domain.PullRequest, *string, error) {
	tx, err := s.database.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		slog.Error("transaction reassign pull request failed", "error", err)
		return nil, nil, err
	}
	defer tx.Rollback(ctx)

	var pr domain.PullRequest
	err = tx.QueryRow(ctx,
		`SELECT id, name, status, author, CreatedAt, MergedAt 
		 FROM pullRequests 
		 WHERE id = $1 FOR UPDATE`,
		pullRequestId).Scan(&pr.PullRequestId, &pr.Name, &pr.Status, &pr.AuthorId, &pr.CreatedAt, &pr.MergedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Info("pull request not found", "pullRequestId", pullRequestId)
			return nil, nil, ErrPullRequestNotFound
		}
		return nil, nil, err
	}

	if pr.Status == "MERGED" {
		return &pr, nil, ErrCannotReassignMergedPullRequest
	}

	tag, err := tx.Exec(ctx, `DELETE FROM pull_request_reviewers WHERE pull_request = $1 AND reviewer = $2`, pullRequestId, oldReviewerId)
	if err != nil {
		return nil, nil, err
	}
	if tag.RowsAffected() == 0 {
		return &pr, nil, ErrPullRequestReviewerNotFound
	}

	var newReviewerID *string
	err = tx.QueryRow(ctx, `
		SELECT user_id 
		FROM users 
		WHERE team_name = (SELECT team_name FROM users WHERE user_id = $1)
		  AND user_id != $2
		  AND user_id NOT IN (SELECT reviewer FROM pull_request_reviewers WHERE pull_request = $3)
		LIMIT 1
	`, pr.AuthorId, oldReviewerId, pullRequestId).Scan(&newReviewerID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, nil, err
	}

	if newReviewerID != nil {
		_, err = tx.Exec(ctx, `INSERT INTO pull_request_reviewers (pull_request, reviewer) VALUES ($1, $2)`, pullRequestId, *newReviewerID)
		if err != nil {
			return nil, nil, err
		}
		pr.AssignedReviewers = append(pr.AssignedReviewers, *newReviewerID)
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error("transaction reassign commit failed", "error", err)
		return nil, nil, err
	}

	return &pr, newReviewerID, nil
}

func NewPullRequestService(database *infrasctucture.PostgresDb) *PullRequestService {
	return &PullRequestService{database}
}
