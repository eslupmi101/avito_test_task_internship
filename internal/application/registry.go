package application

import (
	domain_service "github.com/example/avito_test_task_internship/internal/domain/service"
	"github.com/example/avito_test_task_internship/internal/infrasctucture"
)

var (
	TeamServiceInstance        *domain_service.TeamService
	UserServiceInstance        *domain_service.UserService
	PullRequestServiceInstance *domain_service.PullRequestService
)

func InitRegistry(database *infrasctucture.PostgresDb) {
	TeamServiceInstance = domain_service.NewTeamService(database)
	UserServiceInstance = domain_service.NewUserService(database)
	PullRequestServiceInstance = domain_service.NewPullRequestService(database)
}
