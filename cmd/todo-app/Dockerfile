FROM golang:1.27.1-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/todo-app/exe /app/cmd/todo-app

FROM alpine:3.24
WORKDIR /app
COPY --from=builder /app/cmd/todo-app/exe /app
COPY --from=builder /app/public /app/public
ENV PROJECT_ROOT=/app
CMD ["/app/exe"]
