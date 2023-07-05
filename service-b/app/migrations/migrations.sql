CREATE TABLE account
(
    id      uuid PRIMARY KEY,
    balance float check ( balance >= 0 ) NOT NULL
);