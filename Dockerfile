# Build (the image weights only onl a 18 MB including the static website)
#
#    export DOCKER_BUILDKIT=1 
#    docker  build -t quid .
#    podman  build -t quid .
#    buildah build -t quid .
#
# Run in prod as a deamon (-d)
#
#    docker run --rm -d -p 0.0.0.0:8082:8082 -e DB_PWD=my_password --name quid quid -env
#    podman run --rm -d -p 0.0.0.0:8082:8082 -e DB_PWD=my_password --name quid quid -env
#
# Run in dev. mode with local PostegeSQL
#
#    docker run --rm --network=host --name quid quid -env
#    podman run --rm --network=host --name quid quid -env

# --------------------------------------------------------------------
FROM docker.io/node:17-alpine AS ui_builder

WORKDIR /code

COPY ui/package.json \
     ui/yarn.lock   ./

RUN set -x                          &&\
    node --version                  &&\
    yarn --version                  &&\
    yarn install --frozen-lockfile  &&\
    yarn cache clean

COPY ui/index.html         \
     ui/postcss.config.js  \
     ui/tailwind.config.js \
     ui/tsconfig.json      \
     ui/vite.config.ts    ./

COPY ui/public public
COPY ui/src    src

RUN set -x        &&\
    ls -lA        &&\
    yarn build

# --------------------------------------------------------------------
FROM docker.io/golang:1.18-alpine AS go_builder

WORKDIR /code

COPY go.mod go.sum ./

RUN set -x            &&\
    ls -lA            &&\
    go version        &&\
    go mod download   &&\
    go mod verify

COPY main.go main.go
COPY quidlib quidlib

# Go build flags "-s -w" removes all debug symbols: https://pkg.go.dev/cmd/link
RUN set -x                                                  &&\
    ls -lA                                                  &&\
    export CGO_ENABLED=0                                    &&\
    export GOFLAGS="-trimpath -modcacherw"                  &&\
    export GOLDFLAGS="-d -s -w -extldflags=-static"         &&\
    go build -a -tags osusergo,netgo -installsuffix netgo   &&\
    ls -sh quid                                             &&\
    ./quid -help     # smoke test

# --------------------------------------------------------------------
FROM docker.io/golang:1.18-alpine AS integrator

WORKDIR /target

# HTTPS root certificates (adds about 200 KB)
# Creaate user & group files
RUN set -x                                                  &&\
    mkdir -p                                 etc/ssl/certs  &&\
    cp -a /etc/ssl/certs/ca-certificates.crt etc/ssl/certs  &&\
    echo 'quid:x:6606:6606::/:' > etc/passwd                &&\
    echo 'quid:x:6606:'         > etc/group

# Static website and backend
COPY --from=ui_builder /code/dist ui/dist
COPY --from=go_builder /code/quid .

# --------------------------------------------------------------------
FROM scratch AS final

COPY --chown=6606:6606 --from=integrator /target /

# Run as unprivileged
USER quid:quid

# Arguments that can be set.
# - at build time using flag "--build-arg" (default env. var. value)
# - at run time using flag "--env" (or -e)
ARG DB_USR="pguser"
ARG DB_PWD="my_password"
ARG DB_HOST=""
ARG QUID_KEY="4f10515b3488a2485a32cf68092b66f195c14b86ac89362e8246661bd2c05c3b"
ARG QUID_ADMIN_USER="admin"
ARG QUID_ADMIN_PWD="my_API_administrator_password"

# URL format is postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]
# see https://stackoverflow.com/a/20722229
ENV DATABASE_URL "postgres://${DB_USR}:${DB_PWD}@${DB_HOST}:5432/quid?sslmode=disable"
#-- DATABASE_URL "dbname=quid user=${DB_USR} password=${DB_PWD} sslmode=disable"

# Configuration used to initialize the Database
ENV QUID_KEY        "${QUID_KEY}"
ENV QUID_ADMIN_USER "${QUID_ADMIN_USER}"
ENV QUID_ADMIN_PWD  "${QUID_ADMIN_PWD}"

# Exposed port
ENV PORT 8082
EXPOSE   8082

# Time zone by default = UTC
ENV TZ UTC0

# Executable name
ENTRYPOINT ["/quid"]

# Default argument(s) appended to ENTRYPOINT
CMD ["-env"]
