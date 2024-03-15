ALTER TABLE book
ADD COLUMN author VARCHAR(255);

UPDATE book
SET author = Author.name
FROM Author
WHERE book.id_author = Author.id;

ALTER TABLE book
DROP CONSTRAINT fk_author;

ALTER TABLE book
DROP COLUMN id_author;

DROP TABLE Author;