package db

// Schema : the PostgreSQL schema.
var schema = `CREATE TABLE IF NOT EXISTS namespace (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	alg TEXT NOT NULL,
	access_key BYTEA NOT NULL,
	refresh_key BYTEA NOT NULL,
	max_access_ttl TEXT NOT NULL DEFAULT '20m',
	max_refresh_ttl TEXT NOT NULL DEFAULT '24h',
	public_endpoint_enabled BOOLEAN NOT NULL DEFAULT false
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

CREATE TABLE IF NOT EXISTS orgtable (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	properties JSONB
);

CREATE INDEX IF NOT EXISTS orgtable_name_idx ON orgtable(name);

CREATE TABLE IF NOT EXISTS usertable (
	id SERIAL PRIMARY KEY,
	username TEXT NOT NULL,
	password TEXT,
	namespace_id INTEGER NOT NULL,
	org_id INTEGER,
	date_created DATE NOT NULL DEFAULT CURRENT_DATE,
	is_disabled BOOLEAN DEFAULT false,
	properties JSONB,
	FOREIGN KEY(namespace_id) REFERENCES namespace(id) ON DELETE RESTRICT,
	FOREIGN KEY(org_id) REFERENCES orgtable(id) ON DELETE RESTRICT,
	UNIQUE (username, namespace_id)
);

CREATE INDEX IF NOT EXISTS user_name_idx ON usertable(username);
CREATE INDEX IF NOT EXISTS org_id_idx ON usertable(org_id);

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

CREATE TABLE IF NOT EXISTS userorg (
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL,
	org_id INTEGER NOT NULL,
	FOREIGN KEY(user_id) REFERENCES usertable(id) ON DELETE CASCADE,
	FOREIGN KEY(org_id) REFERENCES orgtable(id) ON DELETE CASCADE,
	UNIQUE (user_id, org_id)
);

CREATE INDEX IF NOT EXISTS userorg_user_idx ON userorg(user_id);
CREATE INDEX IF NOT EXISTS userorg_org_idx ON userorg(org_id);

CREATE TABLE IF NOT EXISTS namespaceadmin (
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL,
	namespace_id INTEGER NOT NULL,
	FOREIGN KEY(user_id) REFERENCES usertable(id) ON DELETE CASCADE,
	FOREIGN KEY(namespace_id) REFERENCES namespace(id) ON DELETE CASCADE,
	UNIQUE (user_id, namespace_id)
);

CREATE INDEX IF NOT EXISTS namespaceadmin_user_idx ON namespaceadmin(user_id);
CREATE INDEX IF NOT EXISTS namespaceadmin_namespace_idx ON namespaceadmin(namespace_id);

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

CREATE INDEX IF NOT EXISTS token_usr_idx ON token(usr_id);
`

// Schema : the PostgreSQL schema.
var dropAll = `

DROP INDEX IF EXISTS namespaces_name_idx;
DROP INDEX IF EXISTS groups_name_idx;
DROP INDEX IF EXISTS organizations_name_idx;
DROP INDEX IF EXISTS organizations_id_idx;
DROP INDEX IF EXISTS users_name_idx;
DROP INDEX IF EXISTS user_groups_usr_idx;
DROP INDEX IF EXISTS user_groups_grp_idx;
DROP INDEX IF EXISTS user_organizations_usr_idx;
DROP INDEX IF EXISTS user_organizations_org_idx;
DROP INDEX IF EXISTS administrators_usr_idx;
DROP INDEX IF EXISTS administrators_ns_idx;
DROP INDEX IF EXISTS token_usr_idx;

DROP INDEX IF EXISTS namespace_name_idx;
DROP INDEX IF EXISTS grouptable_name_idx;
DROP INDEX IF EXISTS orgtable_name_idx;
DROP INDEX IF EXISTS user_name_idx;
DROP INDEX IF EXISTS org_id_idx;
DROP INDEX IF EXISTS usergroup_user_idx;
DROP INDEX IF EXISTS usergroup_group_idx;
DROP INDEX IF EXISTS userorg_user_idx;
DROP INDEX IF EXISTS userorg_org_idx;
DROP INDEX IF EXISTS namespaceadmin_user_idx;
DROP INDEX IF EXISTS namespaceadmin_namespace_idx;
DROP INDEX IF EXISTS token_user_idx;

DROP TABLE IF EXISTS namespace,
                     groups, grouptable,
					 organizations, orgtable,
					 users, usertable,
					 user_groups, usergroup,
					 user_organizations, userorg,
					 administrators, namespaceadmin,
					 token;
`
