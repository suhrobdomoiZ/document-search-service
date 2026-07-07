CREATE TABLE IF NOT EXISTS documents (
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL,
    rubrics TEXT[] NOT NULL,
    created_date TIMESTAMP NOT NULL
);
