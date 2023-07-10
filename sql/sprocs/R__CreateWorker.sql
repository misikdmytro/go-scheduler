CREATE OR REPLACE PROCEDURE create_worker(
        p_id UUID,
        p_name VARCHAR(255),
        p_topic VARCHAR(255)
    ) LANGUAGE SQL AS $$
INSERT INTO workers(id, name, topic)
VALUES (p_id, p_name, p_topic);
$$;