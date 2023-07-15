CREATE TABLE IF NOT EXISTS workers(
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS workers_name_uindex ON workers (name);
CREATE TABLE IF NOT EXISTS jobs(
    id VARCHAR(36) PRIMARY KEY,
    message TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    output JSON
);