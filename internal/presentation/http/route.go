package presentation_http

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	r.Post("/team/add", addTeam)
	r.Get("/team", getTeam)

	r.Post("/users/setIsActive", setIsActive)
	r.Get("/users/getReview", getReview)

	r.Post("/pullRequest/create", create)
	r.Post("/pullRequest/merge", merge)
	r.Post("/pullRequest/reassigne", getPlaylist)
}
