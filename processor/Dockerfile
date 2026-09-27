FROM golang:1.27.0-alpine AS builder

WORKDIR /app

COPY go.mod ./
# COPY go.sum ./

RUN go mod download

COPY src ./

RUN CGO_ENABLED=0 GOOS=linux go build -o exec

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/exec ./

EXPOSE 8000

CMD ["./exec"]