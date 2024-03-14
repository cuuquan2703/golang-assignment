CREATE TABLE IF NOT EXISTS Book (
    isbn VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    publish_year INT,
    author VARCHAR(255)
);

CREATE OR REPLACE PROCEDURE insert_book(
    IN p_isbn TEXT,
    IN p_name TEXT,
    IN p_publish_year INT,
    IN p_author TEXT
)
LANGUAGE SQL
AS $$
    INSERT INTO book (isbn, name, publish_year, author)
    VALUES (p_isbn, p_name, p_publish_year, p_author);
$$;

