-- Удаление индекса
DROP INDEX IF EXISTS idx_pull_request_reviewer;

-- Уникальность: один и тот же reviewer не может быть добавлен дважды к одному PR
ALTER TABLE pull_request_reviewers
DROP CONSTRAINT IF EXISTS unique_reviewer_per_pr;

-- Ограничение: ревьюер не может быть автором PR
DROP TRIGGER IF EXISTS trg_reviewer_not_author ON pull_request_reviewers;
DROP FUNCTION IF EXISTS check_reviewer_not_author();

-- Ограничение: максим 2 ревьюеров на один pull_request
DROP TRIGGER IF EXISTS trg_max_reviewers ON pull_request_reviewers;
DROP FUNCTION IF EXISTS check_max_reviewers();

DROP TABLE IF EXISTS pull_request_reviewers;
