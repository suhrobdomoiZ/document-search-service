CREATE TABLE IF NOT EXISTS documents (
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL,
    rubrics TEXT[] NOT NULL,
    created_date TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_documents_created_date ON documents(created_date DESC);
