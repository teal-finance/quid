# Setup the Postgresql database

### Check PostreSQL

Quid expects PostreSQL listens to the port 5432.

You can check your PostreSQL status and the port:

```bash
$ sudo service postgresql status
$ ss -nlt | grep 5432
LISTEN  0        244            127.0.0.1:5432           0.0.0.0:*
```

### Create user and database

If you do not have already created a priviledged user, create it:

```bash
$ sudo -u postgres psql
postgres=# create user pguser with password 'my_password';
CREATE ROLE
```

Create the Quid database:

```bash
$ sudo -u postgres psql
postgres=# create database quid;
CREATE DATABASE
postgres=# GRANT ALL PRIVILEGES ON DATABASE quid to pguser;
GRANT
```

You may replace the above last statement by:

```bash
postgres=# \c quid
postgres=# GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public to pguser;
```