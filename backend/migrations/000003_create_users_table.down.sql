-- Drop follow_requests table
DROP TABLE IF EXISTS follow_requests;

-- Drop user_followings table
DROP TABLE IF EXISTS user_followings;

-- Drop FK from movie_review_reactions to users
ALTER TABLE movie_review_reactions
    DROP CONSTRAINT IF EXISTS fk_movie_review_reactions_user;

-- Drop FK from movie_reviews to users
ALTER TABLE movie_reviews
    DROP CONSTRAINT IF EXISTS fk_movie_reviews_user;

-- Drop user_id column from movie_reviews
ALTER TABLE movie_reviews
    DROP COLUMN IF EXISTS user_id;

-- Drop users table
DROP TABLE IF EXISTS users;
