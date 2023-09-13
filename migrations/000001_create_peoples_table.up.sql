CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE sex_enum AS ENUM('male', 'female');

CREATE TABLE IF NOT EXISTS peoples
(
    id         uuid DEFAULT uuid_generate_v4(),
    first_name VARCHAR  NOT NULL,
    last_name  VARCHAR  NOT NULL,
    patronymic VARCHAR  NOT NULL,
    age        INT      NOT NULL,
    sex        sex_enum NOT NULL,
    nation     VARCHAR  NOT NULL,
    deleted    BOOLEAN  NOT NULL DEFAULT FALSE,
    PRIMARY KEY (id)
);