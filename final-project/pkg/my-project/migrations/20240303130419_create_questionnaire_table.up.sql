CREATE TABLE IF NOT EXISTS questionnaire
(
    id        bigserial PRIMARY KEY,
    createdAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    topic     text,
    questions JSONB,
    userId    int8
); 