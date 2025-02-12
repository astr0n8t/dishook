# Build Stage
ARG BUILDPLATFORM
FROM --platform=${BUILDPLATFORM} golang:1.24.0 AS build-stage

LABEL app="dishook"
LABEL REPO="https://github.com/astr0n8t/dishook"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
# Copy all internal modules
COPY cmd/*.go ./cmd/
COPY config/*.go ./config/
COPY internal/*.go ./internal/
COPY version/*.go ./version/

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /dishook

# Deploy the application binary into a lean image
FROM gcr.io/distroless/static-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /dishook /dishook

USER nonroot:nonroot

ENTRYPOINT ["/dishook"]
