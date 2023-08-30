CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS user_data_search_index ON user_data USING gin (first_name gin_trgm_ops, second_name gin_trgm_ops);