DROP TRIGGER IF EXISTS trg_set_merged_at ON pull_requests;
DROP FUNCTION IF EXISTS set_merged_at();

DROP TABLE IF EXISTS pull_requests;
