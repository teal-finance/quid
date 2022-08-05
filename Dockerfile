# Build a tiny image of 20 MB including the static website:
#
#    export DOCKER_BUILDKIT=1
#    docker  build -t quid .
#    podman  build -t quid .
#    buildah build -t quid .
#
# Run in prod as a daemon (-d)
#
#    docker run --rm -d -p 0.0.0.0:8082:8082 -e DB_PWD=my_password --name quid quid -env
#    podman run --rm -d -p 0.0.0.0:8082:8082 -e DB_PWD=my_password --name quid quid -env
#
# Run in dev. mode with local PostgreSQL
#
#    docker run --rm --network=host --name quid quid -env
#    podman run --rm --network=host --name quid quid -env

# Arguments with default values to run Quid as unprivileged.
# Set arguments at build time:
# docker build --build-arg UID=1122 --build-arg GID=1122
ARG UID=6606 \
    GID=6606

# --------------------------------------------------------------------
FROM docker.io/node:18-alpine AS ui-builder

WORKDIR /code

COPY ui/package.json \
    ui/yarn.lock    ./

RUN set -ex                         ;\
    node --version                  ;\
    yarn --version                  ;\
    yarn install --frozen-lockfile  ;\
    yarn cache clean

COPY ui/index.html        \
    ui/postcss.config.js  \
    ui/tailwind.config.js \
    ui/tsconfig.json      \
    ui/vite.config.ts    ./

COPY ui/public public
COPY ui/src    src

RUN set -ex     ;\
    ls -lA      ;\
    yarn build

# --------------------------------------------------------------------
FROM docker.io/golang:1.19.0-alpine AS go-builder

WORKDIR /code

COPY go.mod go.sum ./

RUN set -ex           ;\
    ls -lA            ;\
    go version        ;\
    go mod download   ;\
    go mod verify

COPY main.go main.go
COPY quidlib quidlib

# Go build flags "-s -w" removes all debug symbols: https://pkg.go.dev/cmd/link
RUN set -ex                                                ;\
    ls -lA                                                 ;\
    CGO_ENABLED=0                                           \
    GOFLAGS="-trimpath -modcacherw"                         \
    GOLDFLAGS="-d -s -w -extldflags=-static"                \
    go build -a -tags osusergo,netgo -installsuffix netgo  ;\
    ls -sh quid                                            ;\
    ./quid -help  # smoke test

# --------------------------------------------------------------------
FROM docker.io/golang:1.19.0-alpine AS integrator

WORKDIR /target

ARG UID \
    GID

# HTTPS root certificates (adds about 200 KB).
# Create user & group files.
RUN set -ex                                             ;\
    mkdir -p                              etc/ssl/certs ;\
    cp /etc/ssl/certs/ca-certificates.crt etc/ssl/certs ;\
    echo 'quid:x:${UID}:${GID}::/:' > etc/passwd        ;\
    echo 'quid:x:${GID}:'           > etc/group

# Copy the static website and backend executable.
COPY --from=ui-builder /code/dist ui/dist
COPY --from=go-builder /code/quid .

# --------------------------------------------------------------------
FROM scratch AS final

# Run as unprivileged.
ARG UID \
    GID
USER "${UID}:${GID}"

# In this tiny image, put only the the static website,
# the executable "quid", the SSL certificates,
# the "passwd" and "group" files. No shell commands.
COPY --chown="${UID}:${GID}" --from=integrator /target /

ARG DB_USR=pguser       \
    DB_PWD=my_password  \
    DB_HOST=""

# DATABASE_URL format: postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]
# See https://stackoverflow.com/a/20722229
ENV DATABASE_URL="postgres://${DB_USR}:${DB_PWD}@${DB_HOST}:5432/quid?sslmode=disable"

# QUID_ADMIN_* and QUID_KEY are used to initialize the Database.
ARG QUID_ADMIN_USER=admin                         \
    QUID_ADMIN_PWD=my_API_administrator_password  \
    QUID_KEY=4f10515b3488a2485a32cf68092b66f195c14b86ac89362e8246661bd2c05c3b

ENV QUID_ADMIN_USER=$QUID_ADMIN_USER  \
    QUID_ADMIN_PWD=$QUID_ADMIN_PWD    \
    QUID_KEY=$QUID_KEY

# PORT is the web+API port exposed outside of the container.
ENV PORT=8082
EXPOSE ${PORT}

# Default timezone is UTC.
ENV TZ=UTC0

# The default command to run the container.
ENTRYPOINT ["/quid"]

# Default argument(s) appended to ENTRYPOINT.
CMD ["-env"]
