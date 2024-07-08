CREATE TABLE people (
    passport CHARACTER VARYING(11) PRIMARY KEY,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    patronymic TEXT,
    address TEXT
);