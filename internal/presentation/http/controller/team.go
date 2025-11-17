package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/example/avito_test_task_internship/internal/application"
	domain "github.com/example/avito_test_task_internship/internal/domain/entity"
	domain_service "github.com/example/avito_test_task_internship/internal/domain/service"
)

type memberJSON struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type teamJSON struct {
	TeamName string       `json:"team_name"`
	Members  []memberJSON `json:"members"`
}

func usersFromTeam(team teamJSON) []domain.User {
	users := make([]domain.User, 0, len(team.Members))
	for _, m := range team.Members {
		users = append(users, domain.User{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
			TeamName: team.TeamName,
		})
	}
	return users
}

func teamFromUsers(users []*domain.User) teamJSON {
	var t teamJSON
	if len(users) == 0 {
		return t
	}

	// Берём имя команды из первого пользователя
	t.TeamName = users[0].TeamName

	members := make([]memberJSON, 0, len(users))
	for _, u := range users {
		members = append(members, memberJSON{
			UserID:   u.UserID,
			Username: u.Username,
			IsActive: u.IsActive,
		})
	}

	t.Members = members
	return t
}

func AddTeam(w http.ResponseWriter, r *http.Request) {
	var team teamJSON
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	teamUsers, err := application.TeamServiceInstance.CreateTeam(ctx, team.TeamName, usersFromTeam(team))
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		switch {
		case errors.Is(err, domain_service.ErrTeamExists):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]map[string]string{
				"error": {
					"code":    "TEAM_EXISTS",
					"message": "team_name already exists",
				},
			})
			return
		case errors.Is(err, domain_service.ErrTeamMemberAlreadyExists):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]map[string]string{
				"error": {
					"code":    "TEAM_MEMBER_EXISTS",
					"message": "user_id or username already exists",
				},
			})
			return
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]teamJSON{"team": teamFromUsers(teamUsers)})
}

func GetTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		http.Error(w, "team_name query param is required", http.StatusBadRequest)
		return
	}

	users, err := application.TeamServiceInstance.GetTeam(ctx, teamName)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		if errors.Is(err, domain_service.ErrTeamNotFound) {
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

	team := teamFromUsers(users)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]teamJSON{
		"team": team,
	})
}
