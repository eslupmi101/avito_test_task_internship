package application

import (
	domain "github.com/example/avito_test_task_internship/internal/domain/service"
	"github.com/example/avito_test_task_internship/internal/infrasctucture"
)

var (
	TeamServiceInstance        *domain.TeamService
	UserServiceInstance        *domain.UserService
	PullRequestServiceInstance *domain.PullRequestService
)

func InitRegistry(database *infrasctucture.PostgresDb) {
	TeamServiceInstance = domain.NewTeamService(database)
	UserServiceInstance = domain.NewUserService(database)
	PullRequestServiceInstance = domain.NewPullRequestService(database)
}
