CREATE TABLE IF NOT EXISTS workers(
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL
);
CREATE TABLE IF NOT EXISTS jobs(
    id VARCHAR(36) PRIMARY KEY,
    status SMALLINT NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    output JSON
);