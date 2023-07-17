CREATE TABLE IF NOT EXISTS workers(
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS workers_name_uindex ON workers (name);
CREATE TABLE IF NOT EXISTS jobs(
    id BIGSERIAL PRIMARY KEY,
    job_id VARCHAR(36),
    message VARCHAR(500) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    output JSON
);
CREATE INDEX IF NOT EXISTS jobs_job_id_index ON jobs (job_id);