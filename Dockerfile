FROM golang:alpine

WORKDIR /go/src/app

ENV GOFLAGS="-buildvcs=false" 

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go run migrate/migrate.go

CMD ["go", "run", "."]
