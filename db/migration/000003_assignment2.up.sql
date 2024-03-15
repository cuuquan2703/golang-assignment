CREATE TABLE book_author (
    id_book VARCHAR(255),
    id_author INT,
    FOREIGN KEY (id_book) REFERENCES book(isbn),
    FOREIGN KEY (id_author) REFERENCES author(id)
);

INSERT INTO author (name)
SELECT DISTINCT id_author
FROM book;

INSERT INTO book_author (id_book, id_author)
SELECT b.isbn, a.id
FROM book b
JOIN author a ON b.id_author = a.id;

ALTER TABLE book
DROP COLUMN id_author;