#!/bin/bash

# Script to create user, database, permission on first run only.
# doc: https://github.com/docker-library/docs/blob/master/postgres/README.md#initialization-scripts

# user pguser already exist => skip errors

# TODO: replace PASSWORD 'myDBpwd' by "$POSTGRES_DB"

psql -v --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-END
	CREATE USER pguser WITH PASSWORD 'myDBpwd';
	CREATE DATABASE quid;
	GRANT ALL PRIVILEGES ON DATABASE quid TO pguser;
END
