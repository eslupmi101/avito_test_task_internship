-- Удаление индекса
DROP INDEX IF EXISTS idx_pull_request_reviewer;

-- Уникальность: один и тот же reviewer не может быть добавлен дважды к одному PR
ALTER TABLE pullRequestReviewers
DROP CONSTRAINT IF EXISTS unique_reviewer_per_pr;

-- Ограничение: ревьюер не может быть автором PR
DROP TRIGGER IF EXISTS trg_reviewer_not_author ON pullRequestReviewers;
DROP FUNCTION IF EXISTS check_reviewer_not_author();

-- Ограничение: максим 2 ревьюеров на один pull_request
DROP TRIGGER IF EXISTS trg_max_reviewers ON pullRequestReviewers;
DROP FUNCTION IF EXISTS check_max_reviewers();

DROP TABLE IF EXISTS pullRequestReviewers;
