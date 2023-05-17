# FROM debian:sid-slim AS build
# RUN apt update && apt-get install pulseaudio -y
# ENV PULSE_SERVER=docker.for.mac.localhost
# WORKDIR /MyApp/
# COPY . .
# RUN apt update && apt install -y libasound2-dev && go mod tidy && go build -o ./app ./cmd/main.go
# CMD ["./app"]

FROM golang:1.20
WORKDIR /MyApp/
COPY . .
RUN apt update && apt install -y libasound2-dev && go mod tidy && go build -o ./app ./cmd/main.go
CMD ["./app"]