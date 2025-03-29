FROM golang:1.24-alpine AS build

WORKDIR /app

COPY . .

RUN go build -o cpu-stress .

FROM alpine:latest

COPY --from=build /app/cpu-stress /usr/local/bin/cpu-stress

ENTRYPOINT ["cpu-stress"]
