# Build stage
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/index ./index/ && \
    go build -o /app/import ./import/ && \
    go build -o /app/web ./web/

# Final image
FROM alpine:3.21

RUN apk add --no-cache ca-certificates ctags

WORKDIR /app

COPY --from=builder /app/index .
COPY --from=builder /app/import .
COPY --from=builder /app/web .

EXPOSE 8080

CMD ["./web"]
