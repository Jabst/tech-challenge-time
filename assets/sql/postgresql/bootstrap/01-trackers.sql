CREATE TABLE IF NOT EXISTS time_tracker (
    id              SERIAL,
    start           TIMESTAMP NOT NULL,
    end             TIMESTAMP,
    name            TEXT NOT NULL,
    deleted         BOOL DEFAULT 'f',
    version         INT DEFAULT 1,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),

    PRIMARY KEY(id)
)
