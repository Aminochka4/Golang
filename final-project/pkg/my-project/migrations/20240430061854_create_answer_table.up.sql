CREATE TABLE IF NOT EXISTS answer
(
    id               bigserial PRIMARY KEY,
    createdAt        timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updatedAt        timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    questionnaireId  bigint NOT NULL REFERENCES questionnaire (id) ON DELETE CASCADE ,
    answer           text,
    userId           bigint REFERENCES users (id) ON DELETE CASCADE
);


SELECT * FROM answer;