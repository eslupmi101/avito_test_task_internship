package domain

import (
	"context"
	"errors"
	"log/slog"

	domain "github.com/example/avito_test_task_internship/internal/domain/entity"
	"github.com/example/avito_test_task_internship/internal/infrasctucture"
	"github.com/jackc/pgx/v5"
)

type TeamService struct {
	database *infrasctucture.PostgresDb
}

func (s *TeamService) Add(ctx context.Context, teamName string, users []domain.User) ([]*domain.User, error) {
	tx, err := s.database.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable, // максимальная изоляция
	})
	if err != nil {
		slog.Error("transacton Add Team failed", "error", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Блокируем всю таблицу users, чтобы избежать race conditions
	_, err = tx.Exec(ctx, "LOCK TABLE users IN ACCESS EXCLUSIVE MODE")
	if err != nil {
		slog.Error("transacton Add Team failed", "error", err)
		return nil, err
	}

	// 1. Проверяем, есть ли команда с таким именем
	var count int
	err = tx.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE team_name = $1", teamName).Scan(&count)
	if err != nil {
		slog.Error("transacton Add Team failed", "error", err)
		return nil, err
	}
	if count > 0 {
		slog.Info("team already exists:", "teamName", teamName)
		return nil, errors.New("team already exists")
	}

	// 2. Проверяем уникальность пользователей
	for _, u := range users {
		var userCount int
		err := tx.QueryRow(ctx, `
			SELECT COUNT(*) FROM users 
			WHERE user_id = $1 OR username = $2
		`, u.UserID, u.Username).Scan(&userCount)
		if err != nil {
			slog.Error("transacton Add Team failed", "error", err)
			return nil, err
		}
		if userCount > 0 {
			slog.Info(
				"user already exists",
				"user_id", u.UserID,
				"username", u.Username,
			)
			return nil, errors.New("user already exists with ID or username")
		}
	}

	// 3. Вставляем пользователей
	result := make([]*domain.User, 0, len(users))
	for _, u := range users {
		_, err := tx.Exec(ctx, `
			INSERT INTO users (user_id, username, is_active, team_name) 
			VALUES ($1, $2, $3, $4)
		`, u.UserID, u.Username, u.IsActive, teamName)
		if err != nil {
			slog.Error("transacton Add Team failed", "error", err)
			return nil, err
		}
		result = append(result, &domain.User{
			UserID:   u.UserID,
			Username: u.Username,
			IsActive: u.IsActive,
			TeamName: teamName,
		})
	}

	// 4. Commit транзакции
	if err := tx.Commit(ctx); err != nil {
		slog.Error("transacton Add Team failed", "error", err)
		return nil, err
	}

	return result, nil
}

func (s *TeamService) Get(ctx context.Context, teamName string) ([]*domain.User, error) {
	rows, err := s.database.Pool.Query(ctx,
		"SELECT user_id, username, is_active, team_name FROM users WHERE team_name = $1",
		teamName,
	)
	if err != nil {
		slog.Error("query Get Team failed", "error", err)
		return nil, err
	}
	defer rows.Close()

	users := make([]*domain.User, 0)
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.UserID, &u.Username, &u.IsActive, &u.TeamName); err != nil {
			slog.Error("query Get Team failed", "error", err)
			return nil, err
		}
		users = append(users, &u)
	}

	if len(users) == 0 {
		slog.Error("team not found", "teamName", teamName)
		return nil, errors.New("team not found")
	}

	return users, nil
}

func NewTeamService(database *infrasctucture.PostgresDb) *TeamService {
	return &TeamService{database}
}
