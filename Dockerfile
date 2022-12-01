FROM golang:1.19.3-alpine3.16

WORKDIR /api

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build main.go

ENV GIN_MODE=release

CMD ["./main"]
