CREATE OR REPLACE FUNCTION get_worker(worker_id VARCHAR(36)) RETURNS TABLE (
        id VARCHAR(36),
        name VARCHAR(255),
        description VARCHAR(255)
    ) AS $$
SELECT id,
    name,
    description
FROM workers
WHERE id = worker_id
LIMIT 1;
$$ LANGUAGE SQL;