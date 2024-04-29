CREATE TABLE IF NOT EXISTS questionnaire
(
    id        bigserial PRIMARY KEY,
    createdAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    topic     text,
    questions text,
    userId    bigserial REFERENCES users (id)
);

SELECT * FROM questionnaire;