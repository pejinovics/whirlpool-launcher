FROM golang:latest AS build
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o tcpserver main.go

FROM alpine:latest
WORKDIR /app

COPY --from=build /app/tcpserver .

EXPOSE 9000
CMD ["./tcpserver"]