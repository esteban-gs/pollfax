CREATE TYPE vote AS ENUM ('like', 'dislike');

CREATE TABLE sentiments(
  id SERIAL PRIMARY KEY,
  sentiment vote,
  voted_on TIMESTAMP
);
