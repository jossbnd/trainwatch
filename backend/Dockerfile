FROM golang:1.25-alpine AS dev
RUN go install github.com/air-verse/air@latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
CMD ["air"]

FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /trainwatch ./cmd/main.go

FROM alpine:3.21 AS prod
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /trainwatch /trainwatch
EXPOSE 8080
ENTRYPOINT ["/trainwatch"]
