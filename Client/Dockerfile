FROM golang:1.20 AS build
WORKDIR /MyApp/
COPY . .
RUN apt update && apt install -y libasound2-dev && go mod tidy && go build -o ./app ./cmd/main.go

FROM golang:1.20
WORKDIR /MyApp/
COPY --from=build ./MyApp/app .
COPY --from=build ./MyApp/music ./music
COPY --from=build ./MyApp/configs ./configs
RUN apt update && apt install -y libasound2-dev
CMD ["./app"]