DROP TRIGGER IF EXISTS trg_set_merged_at ON pullRequests;
DROP FUNCTION IF EXISTS set_merged_at();

DROP TABLE IF EXISTS pullRequests;
