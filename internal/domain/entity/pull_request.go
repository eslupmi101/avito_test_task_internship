package domain

import "time"

type PullRequest struct {
	PullRequestId     string
	Name              string
	Status            string
	AuthorId          string
	CreatedAt         time.Time
	MergedAt          *time.Time
	AssignedReviewers []string
}
