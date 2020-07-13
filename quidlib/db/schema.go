package db

// Schema : the Postgresql schema
var schema = `CREATE TABLE IF NOT EXISTS namespace (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	key TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS namespace_name_idx ON namespace(name);

CREATE TABLE IF NOT EXISTS grouptable (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	namespace_id INTEGER NOT NULL,
	date_created DATE NOT NULL DEFAULT CURRENT_DATE,
	properties JSONB,
	FOREIGN KEY(namespace_id) REFERENCES namespace(id) ON DELETE RESTRICT,
	UNIQUE (name, namespace_id)
);

CREATE INDEX IF NOT EXISTS grouptable_name_idx ON grouptable(name);

CREATE TABLE IF NOT EXISTS usertable (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	password TEXT,
	namespace_id INTEGER NOT NULL,
	date_created DATE NOT NULL DEFAULT CURRENT_DATE,
	properties JSONB,
	FOREIGN KEY(namespace_id) REFERENCES namespace(id) ON DELETE RESTRICT,
	UNIQUE (name, namespace_id)
);

CREATE INDEX IF NOT EXISTS user_name_idx ON usertable(name);

CREATE TABLE IF NOT EXISTS usergroup (
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL,
	group_id INTEGER NOT NULL,
	FOREIGN KEY(user_id) REFERENCES usertable(id) ON DELETE CASCADE,
	FOREIGN KEY(group_id) REFERENCES grouptable(id) ON DELETE CASCADE,
	UNIQUE (user_id, group_id)
);

CREATE INDEX IF NOT EXISTS usergroup_user_idx ON usergroup(user_id);
CREATE INDEX IF NOT EXISTS usergroup_group_idx ON usergroup(group_id);

CREATE TABLE IF NOT EXISTS token (
	id SERIAL PRIMARY KEY,
	value TEXT NOT NULL,
	user_id INTEGER,
	expiration_date DATE NOT NULL,
	namespace_id INTEGER NOT NULL,
	claims JSONB,
	FOREIGN KEY(user_id) REFERENCES usertable(id) ON DELETE CASCADE,
	FOREIGN KEY(namespace_id) REFERENCES namespace(id) ON DELETE CASCADE,
	UNIQUE (user_id, namespace_id)
);

CREATE INDEX IF NOT EXISTS token_user_idx ON token(user_id);
`
