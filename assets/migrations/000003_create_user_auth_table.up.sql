CREATE TABLE user_auth (
    user_id uuid NOT NULL,
    created TIMESTAMPTZ NOT NULL,
    hashed_password TEXT NOT NULL
);

ALTER TABLE user_auth
    ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES user_data (user_id);