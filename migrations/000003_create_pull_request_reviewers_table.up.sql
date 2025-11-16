CREATE TABLE pull_request_reviewers (
    id SERIAL PRIMARY KEY,
    pull_request VARCHAR(50) NOT NULL REFERENCES pull_requests(id),
    reviewer VARCHAR(50) NOT NULL REFERENCES users(user_id)
);

-- Ограничение: максим 2 ревьюеров на один pull_request
CREATE OR REPLACE FUNCTION check_max_reviewers()
RETURNS trigger AS $$
DECLARE
    cnt INT;
BEGIN
    SELECT COUNT(*) INTO cnt
    FROM pull_request_reviewers
    WHERE pull_request = NEW.pull_request;

    IF cnt >= 2 THEN
        RAISE EXCEPTION 'Pull request % already has 2 reviewers', NEW.pull_request;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_max_reviewers
BEFORE INSERT ON pull_request_reviewers
FOR EACH ROW
EXECUTE FUNCTION check_max_reviewers();

-- Ограничение: ревьюер не может быть автором PR
CREATE OR REPLACE FUNCTION check_reviewer_not_author()
RETURNS trigger AS $$
DECLARE
    pr_author INT;
BEGIN
    SELECT author INTO pr_author FROM pull_requests WHERE id = NEW.pull_request;

    IF NEW.reviewer = pr_author THEN
        RAISE EXCEPTION 'Reviewer cannot be the author of the pull request';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_reviewer_not_author
BEFORE INSERT ON pull_request_reviewers
FOR EACH ROW
EXECUTE FUNCTION check_reviewer_not_author();

-- Уникальность: один и тот же reviewer не может быть добавлен дважды к одному PR
ALTER TABLE pull_request_reviewers
ADD CONSTRAINT unique_reviewer_per_pr UNIQUE(pull_request, reviewer);

-- Мультиколоночный индекс для ускорения поиска по сочетанию pull_request и reviewer
CREATE INDEX idx_pull_request_reviewer
ON pull_request_reviewers(pull_request, reviewer);