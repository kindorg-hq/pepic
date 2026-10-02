# Go with CGO and libvips: the tests and the binary share one environment.
FROM alpine:edge AS build-env

ENV GOOS=linux
ENV CGO_CFLAGS_ALLOW="-Xpreprocessor"

RUN apk add --no-cache go gcc g++ vips-dev git
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# The tests. The golden path builds this stage before the image, on every PR
# and every merge; locally: docker build --target test .
FROM build-env AS test
RUN go vet ./... && go test ./...

FROM build-env AS builder
RUN go build -a -o /build/app -ldflags="-s -w -h" .

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
