CREATE OR REPLACE FUNCTION append_job_status(
        job_id VARCHAR(36),
        message VARCHAR(500),
        created_at TIMESTAMP,
        output JSON
    ) RETURNS VARCHAR(36) AS $$
INSERT INTO jobs(job_id, message, created_at, output)
VALUES (job_id, message, created_at, output)
RETURNING id;
$$ LANGUAGE SQL;