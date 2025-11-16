CREATE TABLE pullRequests (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    status VARCHAR(6) NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN','MERGED')),
    author VARCHAR(50) NOT NULL REFERENCES users(user_id),
    CreatedAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    MergedAt TIMESTAMPTZ
);

-- Функция для автоматического заполнения MergedAt при смене статуса на MERGED
CREATE OR REPLACE FUNCTION set_merged_at()
RETURNS trigger AS $$
BEGIN
    IF NEW.status = 'MERGED' AND (OLD.status IS DISTINCT FROM 'MERGED') THEN
        NEW.MergedAt := NOW();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггер на обновление статуса
CREATE TRIGGER trg_set_merged_at
BEFORE UPDATE ON pullRequests
FOR EACH ROW
EXECUTE FUNCTION set_merged_at();
