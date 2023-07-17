CREATE OR REPLACE FUNCTION append_job_status(
        job_id VARCHAR(36),
        message VARCHAR(500),
        t TIMESTAMP,
        o JSON
    ) RETURNS VARCHAR(36) AS $$
INSERT INTO jobs(job_id, message, timestamp, output)
VALUES (job_id, message, t, o)
RETURNING id;
$$ LANGUAGE SQL;