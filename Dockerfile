FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /note-server ./cmd/server

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /note-server /note-server
COPY --from=builder /app/web /app/web

EXPOSE 8080

CMD ["/note-server"]
