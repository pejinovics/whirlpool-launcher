
FROM golang:latest AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /whirlpool-launcher ./main.go

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /whirlpool-launcher /app/whirlpool-launcher

EXPOSE 8000
ENTRYPOINT ["/app/whirlpool-launcher"]
