#!/bin/bash

# Script to create user, database, permission on first run only.
# doc: https://github.com/docker-library/docs/blob/master/postgres/README.md#initialization-scripts

# user pguser alrady exist => skip errors

psql -v --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-END
	CREATE USER pguser WITH PASSWORD 'my_password';
	CREATE DATABASE quid;
	GRANT ALL PRIVILEGES ON DATABASE quid TO pguser;
END
