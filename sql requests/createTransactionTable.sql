CREATE TABLE IF NOT EXISTS transactions (
    transaction_id SERIAL PRIMARY KEY,
	user_id INT NOT NULL,
	amount INT NOT NULL,
	operation TEXT NOT NULL,
    date TIMESTAMP NOT NULL)