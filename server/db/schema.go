package db

// Schema : the PostgreSQL schema.
var schema = `

CREATE TABLE IF NOT EXISTS namespaces (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	alg TEXT NOT NULL,
	access_key BYTEA NOT NULL,
	refresh_key BYTEA NOT NULL,
	max_access_ttl TEXT NOT NULL DEFAULT '20m',
	max_refresh_ttl TEXT NOT NULL DEFAULT '24h',
	public_endpoint_enabled BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX IF NOT EXISTS namespaces_name_idx ON namespaces(name);

CREATE TABLE IF NOT EXISTS groups (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	ns_id INTEGER NOT NULL,
	date_created DATE NOT NULL DEFAULT CURRENT_DATE,
	properties JSONB,
	FOREIGN KEY(ns_id) REFERENCES namespaces(id) ON DELETE RESTRICT,
	UNIQUE (name, ns_id)
);

CREATE INDEX IF NOT EXISTS groups_name_idx ON groups(name);

CREATE TABLE IF NOT EXISTS organizations (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	properties JSONB
);

CREATE INDEX IF NOT EXISTS organizations_name_idx ON organizations(name);

CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	username TEXT NOT NULL,
	password TEXT,
	ns_id INTEGER NOT NULL,
	org_id INTEGER,
	date_created DATE NOT NULL DEFAULT CURRENT_DATE,
	enabled BOOLEAN DEFAULT true,
	properties JSONB,
	FOREIGN KEY(ns_id) REFERENCES namespaces(id) ON DELETE RESTRICT,
	FOREIGN KEY(org_id) REFERENCES organizations(id) ON DELETE RESTRICT,
	UNIQUE (username, ns_id)
);

CREATE INDEX IF NOT EXISTS organizations_id_idx ON users(org_id);

CREATE INDEX IF NOT EXISTS users_name_idx ON users(username);

CREATE TABLE IF NOT EXISTS user_groups (
	id SERIAL PRIMARY KEY,
	usr_id INTEGER NOT NULL,
	grp_id INTEGER NOT NULL,
	FOREIGN KEY(usr_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY(grp_id) REFERENCES groups(id) ON DELETE CASCADE,
	UNIQUE (usr_id, grp_id)
);

CREATE INDEX IF NOT EXISTS user_groups_usr_idx ON user_groups(usr_id);

CREATE INDEX IF NOT EXISTS user_groups_grp_idx ON user_groups(grp_id);

CREATE TABLE IF NOT EXISTS user_organizations (
	id SERIAL PRIMARY KEY,
	usr_id INTEGER NOT NULL,
	org_id INTEGER NOT NULL,
	FOREIGN KEY(usr_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY(org_id) REFERENCES organizations(id) ON DELETE CASCADE,
	UNIQUE (usr_id, org_id)
);

CREATE INDEX IF NOT EXISTS user_organizations_usr_idx ON user_organizations(usr_id);

CREATE INDEX IF NOT EXISTS user_organizations_org_idx ON user_organizations(org_id);

CREATE TABLE IF NOT EXISTS administrators (
	id SERIAL PRIMARY KEY,
	usr_id INTEGER NOT NULL,
	ns_id INTEGER NOT NULL,
	FOREIGN KEY(usr_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY(ns_id) REFERENCES namespaces(id) ON DELETE CASCADE,
	UNIQUE (usr_id, ns_id)
);

CREATE INDEX IF NOT EXISTS administrators_usr_idx ON administrators(usr_id);

CREATE INDEX IF NOT EXISTS administrators_ns_idx ON administrators(ns_id);

CREATE TABLE IF NOT EXISTS token (
	id SERIAL PRIMARY KEY,
	value TEXT NOT NULL,
	usr_id INTEGER,
	expiration_date DATE NOT NULL,
	ns_id INTEGER NOT NULL,
	claims JSONB,
	FOREIGN KEY(usr_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY(ns_id) REFERENCES namespaces(id) ON DELETE CASCADE,
	UNIQUE (usr_id, ns_id)
);

CREATE INDEX IF NOT EXISTS token_usr_idx ON token(usr_id);
`

// Schema : the PostgreSQL schema.
var dropAll = `
DROP INDEX IF EXISTS administrators_ns_idx,
                     administrators_usr_idx,
                     groups_name_idx,
                     grouptable_name_idx,
                     namespace_name_idx,
                     namespaceadmin_namespace_idx,
                     namespaceadmin_user_idx,
                     namespaces_name_idx,
                     org_id_idx,
                     organizations_id_idx,
                     organizations_name_idx,
                     orgtable_name_idx,
                     token_user_idx,
                     token_usr_idx,
                     user_groups_grp_idx,
                     user_groups_usr_idx,
                     user_name_idx,
                     user_organizations_org_idx,
                     user_organizations_usr_idx,
                     usergroup_group_idx,
                     usergroup_user_idx,
                     userorg_org_idx,
                     userorg_user_idx,
                     users_name_idx;
DROP TABLE IF EXISTS groups, grouptable,
					 organizations, orgtable,
					 users, usertable,
					 user_groups, usergroup,
					 user_organizations, userorg,
					 administrators, namespaceadmin,
					 token,
					 namespaces, namespace;
`
