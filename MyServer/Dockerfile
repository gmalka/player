FROM golang:1.20 AS build
WORKDIR /MyApp/
COPY . .
RUN go mod tidy && go build -o ./app ./cmd/main.go

FROM golang:1.20
WORKDIR /MyApp/
COPY --from=build /MyApp/app .
COPY --from=build /MyApp/configs ./configs
COPY --from=build /MyApp/music ./music
CMD ["./app"]