FROM alpine:edge AS builder

ENV GOOS=linux
ENV CGO_CFLAGS_ALLOW="-Xpreprocessor"

RUN apk add --no-cache go gcc g++ vips-dev git
COPY . /build
WORKDIR /build
RUN go get
RUN go build -a -o /build/app -ldflags="-s -w -h" .

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
