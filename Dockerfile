FROM golang:alpine

WORKDIR /go/src/app

ENV GOFLAGS="-buildvcs=false" 

COPY go.mod go.sum ./

RUN go mod download && go mod verify

RUN go install github.com/githubnemo/CompileDaemon@latest

CMD ["CompileDaemon", "-command=./app"]