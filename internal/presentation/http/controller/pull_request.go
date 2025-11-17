package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/example/avito_test_task_internship/internal/application"
	domain_service "github.com/example/avito_test_task_internship/internal/domain/service"
)

type createPRRequestJSON struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type pullRequestResponseJSON struct {
	PullRequestId     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	var req createPRRequestJSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	pr, err := application.PullRequestServiceInstance.Create(
		ctx,
		req.PullRequestId,
		req.AuthorID,
		req.PullRequestName,
	)
	if err != nil {
		if errors.Is(err, domain_service.ErrUserNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]map[string]string{
				"error": {
					"code":    "NOT_FOUND",
					"message": "resource not found",
				},
			})
		} else if errors.Is(err, domain_service.ErrPullRequestAlreadyExists) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]map[string]string{
				"error": {
					"code":    "PR_EXISTS",
					"message": "PR id already exists",
				},
			})
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	resp := pullRequestResponseJSON{
		PullRequestId:     pr.PullRequestId,
		PullRequestName:   pr.Name,
		AuthorID:          pr.AuthorId,
		Status:            pr.Status,
		AssignedReviewers: pr.AssignedReviewers,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"pr": resp,
	})
}

type mergePRRequest struct {
	PullRequestId string `json:"pull_request_id"`
}

func Merge(w http.ResponseWriter, r *http.Request) {
	var req mergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	pr, err := application.PullRequestServiceInstance.Merge(ctx, req.PullRequestId)
	if err != nil {
		if errors.Is(err, domain_service.ErrPullRequestNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]map[string]string{
				"error": {
					"code":    "NOT_FOUND",
					"message": "resource not found",
				},
			})
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	resp := pullRequestResponseJSON{
		PullRequestId:     pr.PullRequestId,
		PullRequestName:   pr.Name,
		AuthorID:          pr.AuthorId,
		Status:            pr.Status,
		AssignedReviewers: pr.AssignedReviewers,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"pr": resp,
	})
}

type reassignRequestJSON struct {
	PullRequestId string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
}

func Reassign(w http.ResponseWriter, r *http.Request) {
	var req reassignRequestJSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	pr, newReviewerId, err := application.PullRequestServiceInstance.Reassign(ctx, req.PullRequestId, req.OldReviewerID)
	if err != nil {
		if errors.Is(err, domain_service.ErrPullRequestNotFound) || errors.Is(err, domain_service.ErrPullRequestReviewerNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]map[string]string{
				"error": {
					"code":    "NOT_FOUND",
					"message": "resource not found",
				},
			})
		} else if errors.Is(err, domain_service.ErrCannotReassignMergedPullRequest) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]map[string]string{
				"error": {
					"code":    "PR_MERGED",
					"message": "cannot reassign on merged PR",
				},
			})
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	resp := map[string]any{
		"pr": map[string]any{
			"pull_request_id":    pr.PullRequestId,
			"pull_request_name":  pr.Name,
			"author_id":          pr.AuthorId,
			"status":             pr.Status,
			"assigned_reviewers": pr.AssignedReviewers,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"pr":          resp,
		"replaced_by": newReviewerId,
	})
}
