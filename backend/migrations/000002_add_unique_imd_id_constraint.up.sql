ALTER TABLE movie_reviews
ADD CONSTRAINT unique_imdb_id UNIQUE (imdb_id);