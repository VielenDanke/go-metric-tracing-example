FROM golang:latest AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /users ./cmd/service

FROM alpine:latest

WORKDIR /

COPY --from=build /users /users

EXPOSE 8080

ENTRYPOINT ["/users"]