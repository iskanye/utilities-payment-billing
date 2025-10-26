CREATE TABLE IF NOT EXISTS bills
(
    id          SERIAL PRIMARY KEY,
    address     TEXT,
    amount      INTEGER,
    due_date    DATE
);