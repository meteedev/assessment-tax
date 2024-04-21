FROM golang:1.21-alpine as build-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test -v ./...

RUN go build -o ./out/wallet-app .

FROM alpine:3.16.2
COPY --from=build-base /app/out/wallet-app /app/ktaxes-app

CMD ["/app/ktaxes-app"]