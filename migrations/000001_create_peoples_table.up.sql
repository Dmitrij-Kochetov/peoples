CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE sex_enum AS ENUM('male', 'female');

CREATE TABLE IF NOT EXISTS peoples
(
    id         uuid DEFAULT uuid_generate_v4(),
    first_name VARCHAR  NOT NULL,
    last_name  VARCHAR  NOT NULL,
    patronymic VARCHAR,
    age        INT,
    sex        sex_enum,
    nation     VARCHAR,
    deleted    BOOLEAN  NOT NULL DEFAULT FALSE,
    PRIMARY KEY (id)
);