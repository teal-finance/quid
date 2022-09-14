# Setup the PostgreSQL database

## Check PostgreSQL

Install, then check the PostgreSQL status and port:

```sh
$ sudo apt install postgresql
$ sudo service postgresql status
$ ss -nlt | grep 5432
LISTEN 0      244        127.0.0.1:5432      0.0.0.0:*
LISTEN 0      244            [::1]:5432         [::]:*
```

Note: By default, Quid connects to the PostgreSQL server on port 5432.

## Configuration

Quid can be configured with command line options, environnement variables and configuration file.

The following table concerns only the database:

| Command line options | Environnement variables | `config.json`                   | Default value                                                       |
| -------------------- | ----------------------- | ------------------------------- | ------------------------------------------------------------------- |
| `-db-usr`            | `DB_USR`                | `"db_user": "pguser",`          | `pguser`                                                            |
| `-db-pwd`            | `DB_PWD`                | `"db_password": "my_password",` | `my_password`                                                       |
| `-db-host`           | `DB_HOST`               |                                 | `localhost`                                                         |
| `-db-port`           | `DB_PORT`               |                                 | `5432`                                                              |
| `-db-name`           | `DB_NAME`               | `"db_name": "quid",`            | `quid`                                                              |
| `-db-url`            | `DB_URL`                |                                 | `postgres://pguser:my_password@localhost:5432/quid?sslmode=disable` |

When the `-db-url` option or the `DB_URL` env. var. is used,
Quid does not uses the other options, env. vars and the configuration file.
The inverse, when `-db-url` and `DB_URL` are still set to the default value,
Quid rewrites the `-db-url` `DB_URL` using the following formula:

```py
URL = "postgres://{USR}:{PWD}@{HOST}:{PORT}/{NAME}?sslmode=disable"
```

## Setup in one command line

```sql
sudo -u postgres psql -c "CREATE USER pguser WITH PASSWORD 'my_password'" -c "CREATE DATABASE quid" -c "GRANT ALL PRIVILEGES ON DATABASE quid TO pguser"
```

output:

```sql
CREATE ROLE
CREATE DATABASE
GRANT
```

The previous command line perform all the setup operations described in the following chapters.

## Create user and database

If you do not have already created a privileged user, create it:

```sql
$ sudo -u postgres psql
postgres=# CREATE USER pguser WITH PASSWORD 'my_password';
CREATE ROLE
postgres=# exit
```

Update the `config.json` file:

```json
"db_user": "pguser",
"db_password": "my_password",
```

## Create the `quid` database

```sql
$ sudo -u postgres psql
postgres=# CREATE DATABASE quid;
CREATE DATABASE
postgres=# exit
```

## Set the database permissions

```sql
$ sudo -u postgres psql
postgres=# GRANT ALL PRIVILEGES ON DATABASE quid TO pguser;
GRANT
postgres=# exit
```

The previous statement may be replaced by:

```sql
$ sudo -u postgres psql
postgres=# \c quid
You are now connected to database "quid" as user "postgres".
quid-# GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO pguser;
GRANT
quid=# exit
```
