CREATE TABLE IF NOT EXISTS clients (
    id              uuid            PRIMARY KEY DEFAULT gen_random_uuid(),
    login           text            NOT NULL,
    email           text            NOT NULL,
    date_create     timestamptz     NULL DEFAULT now(),
    date_update     timestamptz     DEFAULT NULL,
    date_delete     timestamptz     DEFAULT NULL
);