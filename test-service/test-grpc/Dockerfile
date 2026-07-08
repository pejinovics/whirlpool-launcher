FROM golang:latest AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o grpcserver main.go

FROM alpine:latest
WORKDIR /app

COPY --from=build /app/grpcserver .

EXPOSE 50051
CMD ["./grpcserver"]
