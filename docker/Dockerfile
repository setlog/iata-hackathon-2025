# Stage 0. Build the binary
FROM golang:1.21 as builder

ENV CGO_ENABLED=0
ENV GO111MODULE=on

WORKDIR /go/src/myapp
COPY go.mod .
RUN go mod tidy
COPY . .
RUN  go get ./cmd && go build -C ./cmd -o ../hwb-parser

FROM busybox:1.37

ENV fluenthost fluentbit

RUN mkdir /app
WORKDIR /app

# Copying application and config from previous stage
COPY --from=builder /go/src/myapp/hwb-parser hwb-parser

ENTRYPOINT ["./hwb-parser"]
