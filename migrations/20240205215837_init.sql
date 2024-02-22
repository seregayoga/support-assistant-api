-- +goose Up
CREATE EXTENSION vector;
-- vector(1536) because we use text-embedding-3-small model which has 1536 vector length
-- https://platform.openai.com/docs/guides/embeddings/how-to-get-embeddings
CREATE TABLE support_history (id bigserial PRIMARY KEY, chat text, embedding vector(1536));

-- +goose Down
DROP TABLE support_history;
DROP EXTENSION vector;
