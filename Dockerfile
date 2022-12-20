FROM golang:1.19.3-alpine3.16 AS build

WORKDIR /api

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build main.go


FROM alpine:3.17

WORKDIR /api

COPY --from=build /api/main /api/docs ./
RUN mkdir env
COPY --from=build /api/env/.env ./env

ENV GIN_MODE=release

CMD ["./main"]
