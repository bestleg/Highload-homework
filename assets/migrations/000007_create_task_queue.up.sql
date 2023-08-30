CREATE TABLE IF NOT EXISTS tasks_queue
(
    id         BIGSERIAL               NOT NULL,
    name       TEXT                    NOT NULL,
    payload    JSON                    NOT NULL,
    in_queue   BOOL                    NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    CONSTRAINT tasks_queue_pk
    PRIMARY KEY (id)
    );
