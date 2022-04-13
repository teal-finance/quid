# Build a 18 MB image:
#
#    export DOCKER_BUILDKIT=1 
#    docker  build -t quid .
#    podman  build -t quid .
#    buildah build -t quid .
#
# Run:
#
#    docker run -d --rm -p 0.0.0.0:8082:8082 --name quid quid
#    podman run -d --rm -p 0.0.0.0:8082:8082 --name quid quid

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

# Use UTC time zone by default
ENV TZ UTC0

EXPOSE 8082

ENTRYPOINT ["/quid"]
