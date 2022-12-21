	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);

	INSERT INTO expenses (title, amount, note, tags) VALUES ('hi', 1.4, 'some note', '{"tag1", "tag2"}')
