CREATE TABLE "games" (
  id              VARCHAR(255) PRIMARY KEY NOT NULL,
  title           VARCHAR(255),
  genres          VARCHAR(255),
  company         VARCHAR(255),
  score           INTEGER,
  released_at     TIMESTAMP WITHOUT TIME ZONE
)