package domain

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"errors"

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

	// Проверка, что PR с таким ID еще не существует
	var count int
	err = tx.QueryRow(ctx, `SELECT COUNT(*) FROM pullRequests WHERE id = $1`, pullRequestId).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		slog.Info("pull request already exists with ID", "pullRequestId", pullRequestId)
		return nil, errors.New("pull request already exists with ID")
	}

	// Создаем PR
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
			return nil, errors.New("pull request not found")
		}
		return nil, err
	}

	if pr.Status == "MERGED" {
		slog.Info("pull request already merged", "status", pr.Status)
		return &pr, errors.New("pull request already merged")
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

// Reassign убирает старого ревьюера и возвращает ID старого
func (s *PullRequestService) Reassign(ctx context.Context, pullRequestId string, oldReviewerID string) (*domain.PullRequest, string, error) {
	tx, err := s.database.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		slog.Error("transacton reassing pull request failed", "error", err)
		return nil, "", err
	}
	defer tx.Rollback(ctx)

	var pr domain.PullRequest
	err = tx.QueryRow(ctx, `SELECT id, name, status, author, CreatedAt, MergedAt FROM pullRequests WHERE id = $1 FOR UPDATE`,
		pullRequestId).Scan(&pr.PullRequestId, &pr.Name, &pr.Status, &pr.AuthorId, &pr.CreatedAt, &pr.MergedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			slog.Info("pull request not found:", "pullRequestId", pullRequestId)
			return nil, "", fmt.Errorf("pull request not found")
		}
		return nil, "", err
	}

	// Удаляем старого ревьюера
	tag, err := tx.Exec(ctx, `DELETE FROM pull_request_reviewers WHERE pull_request = $1 AND reviewer = $2`, pullRequestId, oldReviewerID)
	if err != nil {
		return nil, "", err
	}
	if tag.RowsAffected() == 0 {
		slog.Info(
			"reviewer %s not found for pull request %s",
			"oldReviewerID", oldReviewerID, "pullRequestId", pullRequestId,
		)
		return nil, "", fmt.Errorf("reviewer not found for pull request")
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error("transacton reassing pull request failed", "error", err)
		return nil, "", err
	}

	return &pr, oldReviewerID, nil
}

func NewPullRequestService(database *infrasctucture.PostgresDb) *PullRequestService {
	return &PullRequestService{database}
}
