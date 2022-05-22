FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o hello main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/hello .
EXPOSE 8080
CMD ["./hello"]