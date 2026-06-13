FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go.mod download
COPY . .
RUN go build -o mefetch-readme .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/mefetch-readme .
COPY --from=builder /app/static ./static
EXPOSE 8080
CMD ["./mefetch-readme"]