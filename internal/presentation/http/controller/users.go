package controller

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/example/avito_test_task_internship/internal/application"
	domain_service "github.com/example/avito_test_task_internship/internal/domain/service"
)

type userJSON struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

func SetIsActive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		UserId   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	slog.Debug(req.UserId)
	u, err := application.UserServiceInstance.SetIsActive(ctx, req.UserId, req.IsActive)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		if errors.Is(err, domain_service.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]map[string]string{
				"error": {
					"code":    "NOT_FOUND",
					"message": "resource not found",
				},
			})
			return
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user := userJSON{
		UserID:   u.UserID,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]userJSON{"user": user})
}

type prJSON struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

func GetReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := r.URL.Query().Get("user_id")

	prs, err := application.UserServiceInstance.GetReviews(ctx, userId)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		if errors.Is(err, domain_service.ErrNoPullRequestUser) {

			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]map[string]string{
				"error": {
					"code":    "NOT_FOUND",
					"message": "resource not found",
				},
			})
			return
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp := make([]prJSON, 0, len(prs))
	for _, pr := range prs {
		resp = append(resp, prJSON{
			PullRequestID:   pr.PullRequestId,
			PullRequestName: pr.Name,
			AuthorID:        pr.AuthorId,
			Status:          pr.Status,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"user_id":       userId,
		"pull_requests": resp,
	})
}
