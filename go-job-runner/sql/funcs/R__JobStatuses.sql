CREATE OR REPLACE FUNCTION job_statuses(j_id VARCHAR(36)) RETURNS TABLE (
        id BIGINT,
        job_id VARCHAR(36),
        message VARCHAR(500),
        created_at TIMESTAMP,
        output JSON
    ) AS $$
SELECT id,
    job_id,
    message,
    created_at,
    output
FROM jobs
WHERE job_id = j_id
ORDER BY created_at DESC;
$$ LANGUAGE SQL;