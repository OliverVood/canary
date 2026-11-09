CREATE TABLE IF NOT EXISTS clients (
    id              BIGINT          GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    uuid            UUID            DEFAULT gen_random_uuid() UNIQUE,
    login           TEXT            NOT NULL,
    email           TEXT            NOT NULL,
    date_create     TIMESTAMPTZ     NULL DEFAULT now(),
    date_update     TIMESTAMPTZ     DEFAULT NULL,
    date_delete     TIMESTAMPTZ     DEFAULT NULL
);