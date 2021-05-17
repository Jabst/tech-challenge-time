FROM golang:1.14.3 AS build

WORKDIR /src
COPY ./backend .

RUN ls

RUN go mod download
RUN GOOS=linux GO_ARCH=amd64 go build -a -installsuffix cgo -o dist/api cmd/main.go

EXPOSE 8080