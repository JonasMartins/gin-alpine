# Build stage
FROM golang:1.22.3-alpine3.18 AS builder
ARG BINARY_DEV=_main
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 \
  go build -o ./${BINARY_DEV} ./src/services/main/cmd/*.go

RUN chmod +x ./${BINARY_DEV}

## small image
FROM alpine:latest
ARG BINARY_DEV=_main
WORKDIR /app

RUN apk add uuidgen

COPY --from=builder /app/${BINARY_DEV} /app
COPY src/services/main/configs/base-docker-dev.yaml /app/project/src/services/main/configs/base.yaml
COPY src/internal/repository/postgres/migrations /app/project/src/internal/repository/postgres/migrations
COPY wait-for.sh .
COPY start.sh .

CMD ["/app/_main"]
ENTRYPOINT [ "/app/start.sh" ]
# ENTRYPOINT ["tail", "-f", "/dev/null"]
