FROM golang:1.20 AS build
WORKDIR /MyApp/
COPY . .
RUN apt update && apt install -y libasound2-dev && go mod tidy && go build -o ./app ./cmd/main.go

FROM golang:1.20
WORKDIR /MyApp/
COPY --from=builder ./app .
COPY --from=builder ./music .
COPY --from=builder ./configs .
RUN apt update && apt install -y libasound2-dev && go build -o ./app ./cmd/main.go
CMD ["./app"]