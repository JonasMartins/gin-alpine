# Build stage
FROM golang:1.22-alpine3.18 AS builder
ARG BINARY_PROD=_main
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 \
  go build -o ./${BINARY_PROD} ./src/services/main/cmd/*.go

RUN chmod +x ./${BINARY_PROD}

## small image
FROM alpine:latest
ARG BINARY_PROD=_main
WORKDIR /app

RUN apk add uuidgen

COPY --from=builder /app/${BINARY_PROD} /app
COPY src/services/main/configs/base-docker-prod.yaml /app/project/src/services/main/configs/base.yaml
COPY src/internal/repository/postgres/migrations /app/project/src/internal/repository/postgres/migrations
COPY start.sh .

CMD ["/app/_main"]
ENTRYPOINT [ "/app/start.sh" ]
# ENTRYPOINT ["tail", "-f", "/dev/null"]
