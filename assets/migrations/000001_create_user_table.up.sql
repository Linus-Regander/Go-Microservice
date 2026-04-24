CREATE TABLE IF NOT EXISTS user (
    id_u CHAR(36) NOT NULL,
    username VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role ENUM('admin', 'intern'),
    created_at_u TIMESTAMP,
    updated_at_u TIMESTAMP,
    PRIMARY KEY (id_u)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;