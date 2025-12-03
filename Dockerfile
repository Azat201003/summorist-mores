FROM golang:1.25

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download


COPY cmd ./cmd
COPY tests ./tests

RUN ls -a

CMD go run ./cmd

