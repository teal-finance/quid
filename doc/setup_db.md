# Setup the PostreSQL database

## Check PostreSQL

Quid expects PostreSQL listens to the port 5432.

You can check your PostreSQL status and the port:

```sh
$ sudo apt install postgresql
$ sudo service postgresql status
$ ss -nlt | grep 5432
LISTEN 0      244        127.0.0.1:5432      0.0.0.0:*
LISTEN 0      244            [::1]:5432         [::]:*
```

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

## Perform all the previous operations in one command line

```sql
sudo -u postgres psql -c "CREATE USER pguser WITH PASSWORD 'my_password'" -c "CREATE DATABASE quid" -c "GRANT ALL PRIVILEGES ON DATABASE quid TO pguser"
CREATE ROLE
CREATE DATABASE
GRANT
```
