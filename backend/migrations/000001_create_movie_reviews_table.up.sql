-- Create enum for reaction types
CREATE TYPE movie_review_reaction AS ENUM (
  'Agree',
  'Insightful',
  'Funny',
  'ThoughtProvoking',
  'Disagree',
  'WellSaid'
);

-- Main movie_reviews table
CREATE TABLE movie_reviews (
                               id BIGSERIAL PRIMARY KEY,
                               imdb_id VARCHAR(20) NOT NULL,
                               rating SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
                               statement_comment TEXT NOT NULL CHECK (char_length(statement_comment) <= 280),
                               statement_created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                               statement_updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                               created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                               updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
                               version BIGINT NOT NULL DEFAULT 1
);

-- Table to store individual reactions
CREATE TABLE movie_review_reactions (
                                        movie_review_id BIGINT NOT NULL REFERENCES movie_reviews(id) ON DELETE CASCADE,
                                        reaction_type movie_review_reaction NOT NULL,
                                        user_id BIGINT NOT NULL,
                                        PRIMARY KEY (movie_review_id, reaction_type, user_id)
);