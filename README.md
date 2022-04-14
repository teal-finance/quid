# Quid

Quid is a JSON Web Tokens (JWT) server.

## Install

Download the latest [release](https://github.com/teal-finance/quid/releases) to run a binary or clone the repository to compile from source.

## Build from source

    make all -j

## Configure

1. Create the default config file:

       ./quid -conf

2. Create the `quid` database: [instructions](doc/setup_db.md)

3. Edit the configuration file to set your PostgreSQL credentials:

        vim config.json

4. Initialize the `quid` database and create the administrator user:

       ./quid -init

    These registered administrator username and password will be required to login the Administration UI.

## Run the backend

    ./quid

See also: [run in dev mode](doc/dev_mode.md)

Quid serves the static web site. Open <http://localhost:8082> to login into the admin interface:

    xdg-open http://localhost:8082

![Screenshot](doc/img/screenshot.png)

## Deploy on Heroku

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/teal-finance/quid)

## Request tokens

Request a refresh token and use it to request access tokens.

### Refresh token

A public endpoint is available to request refresh tokens for namespaces.
A time to live must be provided.

Example: request a refresh token with a 10 minutes lifetime `/token/refresh/10m`

```php
curl localhost:8082/token/refresh/10m          \
     -H 'Content-Type: application/json'       \
     -d '{"namespace":"my_namespace","username":"my_username","password":"my_password"}'
```

Response:

```json
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IzpXVCJ9..."}
```

### Access token

A public endpoint is available to request access tokens for namespaces.
A time to live must be provided.

Example: request an access token with a 10 minutes lifetime `/token/access/10m`

```php
curl localhost:8082/token/access/10m           \
     -H 'Content-Type: application/json'                      \
     -d '{"namespace":"my_namespace","refresh_token":"zpXVCJ9..."}'
```

Response:

```json
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IzpXVCJ9..."}
```

Note: if the requested duration exceeds the max authorized tokens time to live for the namespace the demand will be rejected

## Decode tokens

### Python

```python
import jwt

try:
    payload = jwt.decode(token, key, algorithms=['HS256'])
except jwt.ExpiredSignatureError:
    # ...
```

Example payload:

```yaml
{
    'username': 'my_username', 
    'groups': ['my_group1', 'my_group2'], 
    'orgs': ['org1', 'org2']
    'exp': 1595950745
}
```

`exp` is the expiration timestamp in [Unix time](https://en.wikipedia.org/wiki/Unix_time) format (seconds since 1970).

### Examples

See the [examples](https://github.com/synw/quid_examples) for various backends.

## Client libraries

Client libraries transparently manage the requests to api servers.
If a server returns a 401 Unauthorized response when an access token is expired,
the client library will request a new access token from a Quid server,
using a refresh token, and will retry the request with the new access token.

### Javascript

[Quidjs](https://github.com/teal-finance/quidjs) : the javascript requests library.
