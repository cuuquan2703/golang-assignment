ALTER TABLE book
ADD COLUMN id_author INT;

UPDATE book
SET id_author = a.id
FROM book_author ba
JOIN author a ON ba.id_author = a.id
WHERE ba.id_book = book.isbn;


ALTER TABLE book
ADD CONSTRAINT fk_author
FOREIGN KEY (id_author)
REFERENCES Author(id);

DROP TABLE book_author;
