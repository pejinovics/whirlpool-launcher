FROM golang:latest AS build
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o httpserver main.go

FROM alpine:latest
WORKDIR /app

COPY --from=build /app/httpserver .

EXPOSE 8080
CMD ["./httpserver"]