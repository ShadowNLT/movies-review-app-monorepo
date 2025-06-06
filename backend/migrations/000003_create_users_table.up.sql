-- Create users table
CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,
                       email citext NOT NULL UNIQUE,
                       password_hash bytea NOT NULL,
                       handle VARCHAR(30) NOT NULL UNIQUE,
                       location VARCHAR(100),
                       date_of_birth DATE NOT NULL,
                       is_protected BOOLEAN NOT NULL DEFAULT FALSE,
                       is_activated BOOLEAN NOT NULL DEFAULT FALSE,
                       created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                       version integer NOT NULL DEFAULT 1
);

-- Add user_id to movie_reviews table to track review authors
ALTER TABLE movie_reviews
    ADD COLUMN user_id BIGINT NOT NULL;

ALTER TABLE movie_reviews
    ADD CONSTRAINT fk_movie_reviews_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Add FK constraint for user_id on movie_review_reactions table
ALTER TABLE movie_review_reactions
    ADD CONSTRAINT fk_movie_review_reactions_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Create user_followings table (A follows B)
CREATE TABLE user_followings (
                                 follower_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                 following_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                 created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                                 PRIMARY KEY (follower_id, following_id)
);

-- Create follow_requests table (pending requests for protected accounts)
CREATE TABLE follow_requests (
                                 requester_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                 target_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                 created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                                 PRIMARY KEY (requester_id, target_id)
);
