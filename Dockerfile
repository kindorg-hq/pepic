# Go with CGO and libvips: the tests and the binary share one environment.
# Stage order is for the Build cache (kindorg-hq/ci README, "How Build
# caches"): what changes least comes first.
FROM alpine:edge AS build-env

ENV GOOS=linux
ENV CGO_CFLAGS_ALLOW="-Xpreprocessor"

RUN apk add --no-cache go gcc g++ vips-dev git
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

# The packages pepic imports (tests included), listed from the source; the
# list alone feeds the next stage, so a change to pepic's own code keeps it.
FROM build-env AS deps
COPY . .
RUN go list -deps -test -f '{{if and .Module (not .Module.Main)}}{{.ImportPath}}{{end}}' ./... \
    | sort -u > /deps.txt

# Dependencies compiled once, into Go's build cache in this layer: rebuilt
# only when the list or go.sum changes, so tests and the binary compile
# pepic's own packages only.
FROM build-env AS src
COPY --from=deps /deps.txt /deps.txt
RUN xargs go build < /deps.txt
COPY . .

# The tests. ci v4 (Build) runs this stage before the image, on every PR
# and every merge; locally: docker build --target test .
FROM src AS test
RUN go vet ./... && go test ./...

FROM src AS builder
RUN go build -o /build/app -ldflags="-s -w -h" .

# The image: keep it the last stage, so a plain build yields it.
FROM alpine:latest

# vips-heif: HEIC/HEIF loader (iPhone photos), converted to JPEG on upload
RUN apk --no-cache add ca-certificates mailcap ffmpeg vips vips-heif
COPY --from=builder /build/app /app/pepic
COPY html /app/html
COPY static /app/static
COPY etc/pepic /etc/pepic
WORKDIR /app

EXPOSE 8118

ENTRYPOINT ["/app/pepic", "serve", "--config", "/etc/pepic/config.yml"]
