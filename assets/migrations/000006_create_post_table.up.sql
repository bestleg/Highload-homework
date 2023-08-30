CREATE TABLE user_posts (
                             post_id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
                             user_id uuid NOT NULL,
                             post text NOT NULL,
                             created TIMESTAMPTZ NOT NULL
);

ALTER TABLE user_posts
    ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES user_data (user_id);