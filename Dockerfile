FROM golang:1.19 AS build

# `boilerplate` should be replaced with your project name
WORKDIR /go/src/detectr

COPY . .

RUN go mod download

# Builds the application as a staticly linked one, to allow it to run on alpine
RUN cd cmd/detectr/ && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o detectr .

FROM alpine:latest

WORKDIR /app

RUN mkdir ./bin
COPY ./bin ./bin

COPY --from=build /go/src/detectr/cmd/detectr .

EXPOSE 3000

CMD ["./detectr"]
