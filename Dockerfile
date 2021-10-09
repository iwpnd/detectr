FROM golang:1.17 AS build

# `boilerplate` should be replaced with your project name
WORKDIR /go/src/detectr

COPY . .

RUN go mod download

# Builds the application as a staticly linked one, to allow it to run on alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .

FROM alpine:latest

WORKDIR /app

RUN mkdir ./bin
COPY ./bin ./bin

COPY --from=build /go/src/detectr/app .

EXPOSE 3000

CMD ["./app"]
