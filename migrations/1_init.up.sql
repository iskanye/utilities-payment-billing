CREATE TABLE IF NOT EXISTS bills
(
    id       SERIAL PRIMARY KEY,
    address  TEXT,
    amount   INTEGER,
    user_id  INTEGER,
    due_date DATE
);