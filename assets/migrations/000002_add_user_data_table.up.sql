CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE user_data (
                       user_id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
                       first_name VARCHAR(50),
                       second_name VARCHAR(50),
                       birthdate DATE,
                       sex bool,
                       biography TEXT,
                       city VARCHAR(50)
);

COMMENT ON COLUMN user_data.sex IS 'true = муж/ false = жен';