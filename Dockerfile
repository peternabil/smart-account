# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
COPY .env .
RUN go build -o main main.go


FROM alpine:3.18
WORKDIR /app
COPY .env .
COPY --from=builder /app/main .

EXPOSE 8080
CMD [ "/app/main" ]
