package presentation_http

import "net/http"

type Handlers struct{}

func addTeam(w http.ResponseWriter, r *http.Request) {}

func getTeam(w http.ResponseWriter, r *http.Request) {}

func setIsActive(w http.ResponseWriter, r *http.Request) {}

func getReview(w http.ResponseWriter, r *http.Request) {}

func create(w http.ResponseWriter, r *http.Request) {}

func merge(w http.ResponseWriter, r *http.Request) {}

func getPlaylist(w http.ResponseWriter, r *http.Request) {}
