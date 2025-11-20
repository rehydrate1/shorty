FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o shorty ./cmd/shorty/main.go


FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/shorty .
COPY .env.example .

EXPOSE 8080
CMD [ "./shorty" ]