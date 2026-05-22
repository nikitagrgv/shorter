CREATE TABLE IF NOT EXISTS links
(
    id                BIGINT      NOT NULL,
    short             TEXT        NOT NULL,
    long_url          TEXT        NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL,
    access_token_hash TEXT        NOT NULL,

    CONSTRAINT links_pkey
        PRIMARY KEY (id),

    CONSTRAINT links_short_unique
        UNIQUE (short),

    CONSTRAINT links_access_token_hash_unique
        UNIQUE (access_token_hash)
);
