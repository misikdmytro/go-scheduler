CREATE OR REPLACE PROCEDURE create_worker(
        p_id UUID,
        p_name VARCHAR(255),
        p_description VARCHAR(255)
    ) LANGUAGE SQL AS $$
INSERT INTO workers(id, name, description)
VALUES (p_id, p_name, p_description);
$$;