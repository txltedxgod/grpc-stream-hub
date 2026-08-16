FROM golang:1.22-alpine AS builder

WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download || true

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/grpc-hub ./server

FROM alpine:3.19
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/grpc-hub /usr/local/bin/grpc-hub

EXPOSE 50051
ENTRYPOINT ["grpc-hub"]
CMD ["-port=50051"]
