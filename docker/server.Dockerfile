FROM golang:1.24.4-alpine3.22 as build_stage

WORKDIR /app/build

COPY go/server .
RUN go build -o /app/server.out server.go

FROM alpine:3.22

WORKDIR /app

COPY --from=build_stage /app/server.out ./server.out

CMD [ "/app/server.out" ]