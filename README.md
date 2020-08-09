# Quid

A Json web tokens server

## Install and run

Download the latest [release](release) to run a binary or clone the repository to compile from source

### Create the database

Create a Postgresql database:

   ```bash
   sudo su postgres
   psql
   create database quid;
   \c quid;
   GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public to pguser;
   \q
   exit
   ```

Replace `pguser` by your database user

### Configure

Create a config file:

   ```bash
   ./quid -conf
   ```

Edit the config file to provide your database credentials. Initialize the database and create an admin user:

   ```bash
   ./quid -init
   ```

### Run

   ```bash
   ./quid
   ```

Go to `localhost:8082` to login into the admin interface

![Screenshot](doc/img/screenshot.png)

### Compile from source

   ```bash
   cd quidui
   npm install
   npm run build
   cd ..
   go build
   ```

[Run in dev mode](doc/dev_mode.md)

## Request tokens

Request a refresh token and use it to request access tokens

### Refresh token

A public endpoint is available to request refresh tokens for namespaces. A time to live must be provided. 
Ex: request a refresh token with a 24h lifetime `/token/refresh/24h`:

   ```bash
   curl -d '{"namespace":"my_namespace","username":"my_username","password":"my_password"}' -H \
   "Content-Type: application/json" -X POST http://localhost:8082/token/refresh/24h
   ```

   Response:

   ```bash
   {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IzpXVCJ9..."}
   ```

### Access token

A public endpoint is available to request access tokens for namespaces. A time to live must be provided. 
Ex: request an access token with a 10 minutes lifetime `/token/access/10m`:

   ```bash
   curl -d '{"namespace":"my_namespace","refresh_token":"zpXVCJ9..."}' -H \
   "Content-Type: application/json" -X POST http://localhost:8082/token/access/10m
   ```

Response:

   ```bash
   {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IzpXVCJ9..."}
   ```

Note: if the requested duration exceeds the max authorized tokens time to live for the namespace the demand will be rejected

## Decode tokens

In python:

   ```python
   payload = jwt.decode(token, key, algorithms=['HS256'])
   ```

Example payload:

   ```javascript
   {
      'namespace': 'my_namespace1', 
      'name': 'my_username', 
      'groups': ['my_group1', 'my_group2'], 
      'exp': 1595950745
   }
   ```

`exp` is the expiration timestamp

Check the [python example](example/python)

## Client library

[Javascript](quidui/src/quidjs/requests.js) client library: [example](quidui/src/api.js) usage
