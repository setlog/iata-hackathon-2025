# Stage 0. Build the binary
FROM golang:1.21 as builder

ARG TOKEN
ENV CGO_ENABLED=0
ENV GO111MODULE=on

# We need this to clone private repo's from our gitlab
RUN if [ -z $TOKEN ]; then echo "TOKEN is not set. Skip git configuration."; else git config --global http.sslVerify false; git config --global url."https://gitlab-ci-token:${TOKEN}@gitlab.setlog.lan/".insteadOf "https://gitlab.setlog.lan/"; fi

WORKDIR /go/src/myapp
COPY go.mod .
RUN go mod tidy
COPY . .
RUN go get && go build -o qc-result-validator

# Stage 1. Manage certs
FROM alpine:3 as certs

RUN apk update && apk add ca-certificates && apk add coreutils && rm -rf /var/cache/apk/*
COPY docker/rds-* /usr/local/share/ca-certificates/
RUN update-ca-certificates

# Stage 2. Run the previous-built binary
FROM busybox:1.35

ENV PORT 8800
ENV fluenthost fluentbit

RUN mkdir /app
WORKDIR /app

# certificates to interact with other services
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copying application and config from previous stage
COPY --from=builder /go/src/myapp/qc-result-validator qc-result-validator

RUN addgroup -g 2001 osca && adduser -D -H -u 2001 -G osca osca

EXPOSE $PORT

USER osca

ENTRYPOINT ["./qc-result-validator"]
