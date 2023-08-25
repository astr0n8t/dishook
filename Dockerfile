# Build Stage
FROM golang:1.21 AS build-stage

LABEL app="build-dishook"
LABEL REPO="https://github.com/astr0n8t/dishook"

ENV PROJPATH=/go/src/github.com/astr0n8t/dishook

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/astr0n8t/dishook
WORKDIR /go/src/github.com/astr0n8t/dishook

RUN make build-alpine

# Final Stage
FROM ghcr.io/astr0n8t/dishook:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/astr0n8t/dishook"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/dishook/bin

WORKDIR /opt/dishook/bin

COPY --from=build-stage /go/src/github.com/astr0n8t/dishook/bin/dishook /opt/dishook/bin/
RUN chmod +x /opt/dishook/bin/dishook

# Create appuser
RUN adduser -D -g '' dishook
USER dishook

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/dishook/bin/dishook"]
