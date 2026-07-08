# build stage (isti kao pre)
FROM golang:latest AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /whirlpool-launcher ./main.go

# runtime stage - koristi sliku koja sadrži docker CLI
FROM docker:24-cli

WORKDIR /app
# kopiraj binarni iz build stage-a
COPY --from=builder /whirlpool-launcher /app/whirlpool-launcher

# opcionalno: certs (docker image obično sadrži)
# RUN apk add --no-cache ca-certificates

EXPOSE 8000
EXPOSE 9090
ENTRYPOINT ["/app/whirlpool-launcher"]
