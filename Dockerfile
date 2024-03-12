# Build a tiny image of 20 MB including the static website:
#
#    export DOCKER_BUILDKIT=1
#    docker  build -t quid .
#    podman  build -t quid .
#    buildah build -t quid .
#
# Run in prod as a daemon (-d)
#
#    docker run --rm -d -p 0.0.0.0:8082:8082 -e POSTGRES_PASSWORD=myDBpwd --name quid quid -env
#    podman run --rm -d -p 0.0.0.0:8082:8082 -e POSTGRES_PASSWORD=myDBpwd --name quid quid -env
#
# Run in dev. mode with local PostgreSQL
#
#    docker run --rm --network=host --name quid quid -dev
#    podman run --rm --network=host --name quid quid -dev

# Arguments with default values to run Quid as unprivileged.
#
# Set arguments at build time:
#
#    docker build --build-arg UID=1122 --build-arg GID=0 .
#
ARG UID=6606
ARG GID=6606

# --------------------------------------------------------------------
FROM docker.io/node:19-alpine AS ui-builder

WORKDIR /code

COPY ui/package.json \
     ui/yarn.lock   ./

RUN set -ex                         ;\
    node --version                  ;\
    yarn --version                  ;\
    yarn install --frozen-lockfile  ;\
    yarn cache clean

COPY ui/index.html         \
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
FROM docker.io/golang:1.22-alpine AS go-builder

WORKDIR /code

COPY go.mod go.sum ./

RUN set -ex           ;\
    ls -lShA          ;\
    go version        ;\
    go mod download   ;\
    go mod verify

COPY cmd    cmd
COPY crypt  crypt
COPY server server
COPY tokens tokens

# Go build flags "-s -w" removes all debug symbols: https://pkg.go.dev/cmd/link
# GOAMD64=v3 --> https://github.com/golang/go/wiki/MinimumRequirements#amd64
RUN set -ex                                                          ;\
    ls -lShA                                                         ;\
    CGO_ENABLED=0                                                     \
    GOFLAGS="-trimpath -modcacherw"                                   \
    GOLDFLAGS="-d -s -w -extldflags=-static"                          \
    GOAMD64=v3                                                        \
    GOEXPERIMENT=newinliner                                           \
    go build -a -tags osusergo,netgo -installsuffix netgo ./cmd/quid ;\
    ls -sh quid                                                      ;\
    ./quid -help  # smoke test

# --------------------------------------------------------------------
FROM docker.io/golang:1.22-alpine AS integrator

WORKDIR /target

ARG UID
ARG GID

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
ARG     UID    GID
USER "${UID}:${GID}"

# In this tiny image, put only the the static website,
# the executable "quid", the SSL certificates,
# the "passwd" and "group" files. No shell commands.
COPY --chown="${UID}:${GID}" --from=integrator /target /

# QUID_ADMIN_* and QUID_KEY are used to initialize the Database.
ARG QUID_ADMIN_USR=quid-admin
ARG QUID_ADMIN_PWD=quid-admin-password
ARG QUID_KEY=95c14b86ac89362e8246661bd2c05c3b
ARG POSTGRES_USER=pguser
ARG POSTGRES_PASSWORD=myDBpwd
ARG POSTGRES_DB=quid
ARG DB_HOST=db
ARG DB_PORT=5432
ARG DB_URL

# Default timezone is UTC.
ENV TZ=UTC0
ENV QUID_ADMIN_USR=$QUID_ADMIN_USR
ENV QUID_ADMIN_PWD=$QUID_ADMIN_PWD
ENV QUID_KEY=$QUID_KEY
ENV POSTGRES_USER=$POSTGRES_USER
ENV POSTGRES_PASSWORD=$POSTGRES_PASSWORD
ENV POSTGRES_DB=$POSTGRES_DB
ENV PORT=8082
ENV DB_HOST=$DB_HOST
ENV DB_PORT=$DB_PORT
ENV DB_URL=$DB_URL

# PORT is the web+API port exposed outside of the container.
EXPOSE ${PORT}

# The default command to run the container.
ENTRYPOINT ["/quid"]

# Default argument(s) appended to ENTRYPOINT.
CMD [""]
