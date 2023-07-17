CREATE OR REPLACE FUNCTION create_worker(
        name VARCHAR(255),
        description VARCHAR(255)
    ) RETURNS VARCHAR(36) AS $$
INSERT INTO workers(id, name, description)
VALUES (gen_random_uuid(), name, description)
RETURNING id;
$$ LANGUAGE SQL;