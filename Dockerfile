FROM golang:1.16-alpine as builder
WORKDIR /app
ENV GO111MODULE=on

COPY . .
RUN go mod download
RUN go build -o recipe-count cmd/main.go


FROM alpine:latest
LABEL maintainer="Alexander Kononykhin <a.kononykhin@yandex.ru>"
WORKDIR /root/
COPY --from=builder /app/recipe-count .
COPY --from=builder /app/config.json .
CMD ["./recipe-count"]
