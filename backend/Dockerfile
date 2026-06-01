FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o backend .

FROM alpine:3.23
WORKDIR /app
RUN apk add --no-cache wget
COPY --from=builder /app/backend .
EXPOSE 8080
HEALTHCHECK --interval=5s --timeout=3s --retries=5 \
  CMD wget -qO- http://localhost:8080/health || exit 1
CMD ["./backend"]
