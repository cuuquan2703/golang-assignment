ALTER TABLE book
ADD COLUMN author VARCHAR(255);

UPDATE book
SET author = Author.name
FROM Author
WHERE book.author_id = Author.id;

DROP TABLE Author;

ALTER TABLE book
DROP CONSTRAINT fk_author;