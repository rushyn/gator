-- +goose Up
CREATE TABLE posts(
id UUID PRIMARY KEY,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
published_at TIMESTAMP NOT NULL,
title TEXT NOT NULL,
url TEXT NOT NULL UNIQUE,
desription TEXT NOT NULL,
feed_id UUID NOT NULL REFERENCES feeds
                            ON DELETE CASCADE 
                            ON UPDATE CASCADE,
FOREIGN KEY (feed_id)
REFERENCES feeds(id)
);

-- +goose Down
DROP TABLE posts;



