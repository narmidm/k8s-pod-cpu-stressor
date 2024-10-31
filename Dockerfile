FROM golang:1.19.2-alpine AS build

WORKDIR /app

COPY . .

RUN go build -o cpu-stress .

FROM alpine:latest

COPY --from=build /app/cpu-stress /usr/local/bin/cpu-stress

LABEL maintainer="narmidm"
LABEL version="1.0.0"
LABEL description="A tool to simulate CPU stress on Kubernetes pods."

ENTRYPOINT ["cpu-stress"]
