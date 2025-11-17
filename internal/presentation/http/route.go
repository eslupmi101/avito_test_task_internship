package presentation

import (
	"github.com/example/avito_test_task_internship/internal/presentation/http/controller"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	r.Post("/team/add", controller.AddTeam)
	r.Get("/team/{team_name}", controller.GetTeam)

	r.Post("/users/setIsActive", controller.SetIsActive)
	r.Get("/users/getReview/{user_id}", controller.GetReview)

	r.Post("/pullRequest/create", controller.Create)
	r.Post("/pullRequest/merge", controller.Merge)
	r.Post("/pullRequest/reassigne", controller.Reassign)
}
