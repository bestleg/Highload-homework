CREATE TABLE user_friend (
                           user_id uuid NOT NULL,
                           user_friend_id uuid NOT NULL,
                           created TIMESTAMPTZ NOT NULL
);

ALTER TABLE user_friend
    ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES user_data (user_id);

ALTER TABLE user_friend
    ADD CONSTRAINT unique_user_id_friend_user_id UNIQUE (user_id, user_friend_id);

CREATE INDEX friend_index ON user_friend (user_friend_id);