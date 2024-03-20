ALTER TABLE book
ADD COLUMN id_author SERIAL;

CREATE TABLE Author (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    birth_date VARCHAR(100)
);

INSERT INTO Author (name)
SELECT DISTINCT author
FROM book;

ALTER TABLE book
DROP COLUMN author;

ALTER TABLE book
ADD CONSTRAINT fk_author
FOREIGN KEY (id_author)
REFERENCES Author(id);