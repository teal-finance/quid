# Setup the PostreSQL database

## Check PostreSQL

Quid expects PostreSQL listens to the port 5432.

You can check your PostreSQL status and the port:

```sh
$ sudo apt install postgresqls
$ sudo service postgresql status
$ ss -nlt | grep 5432
LISTEN 0      244        127.0.0.1:5432      0.0.0.0:*
LISTEN 0      244            [::1]:5432         [::]:*
```

## Create user and database

If you do not have already created a privileged user, create it:

```sh
$ sudo -u postgres psql
postgres=# create user pguser with password 'my_password';
CREATE ROLE
```

(enter [CTRL]+[D] to exist)

Update the `config.json` file:

```json
"db_password": "my_password",
"db_user": "pguser",
```

## Create the Quid database

```sh
$ sudo -u postgres psql
postgres=# create database quid;
CREATE DATABASE
postgres=# GRANT ALL PRIVILEGES ON DATABASE quid to pguser;
GRANT
```

You may replace the above last statement by:

```sh
postgres=# \c quid
postgres=# GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public to pguser;
```
