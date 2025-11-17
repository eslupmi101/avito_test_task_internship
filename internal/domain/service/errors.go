package domain_service

import "errors"

var ErrTeamExists = errors.New("team already exists")

var ErrTeamMemberAlreadyExists = errors.New("member-user already exists with ID or username")

var ErrTeamNotFound = errors.New("team not found")

var ErrUserNotFound = errors.New("user not found")

var ErrNoPullRequestUser = errors.New("no pull requests found for user")

var ErrPullRequestAlreadyExists = errors.New("pull request already exists with ID")

var ErrPullRequestNotFound = errors.New("pull request not found")

var ErrPullRequestAlreadyMerged = errors.New("pull request already merged")

var ErrPullRequestReviewerNotFound = errors.New("reviewer not found for pull request")

var ErrCannotReassignMergedPullRequest = errors.New("cannot reasign merged pull request")
