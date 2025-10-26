FROM golang:1.24 AS base

RUN apt-get update \
  && apt-get install -y make openssh-client ca-certificates unzip \
    && update-ca-certificates \
      && apt-get install bash

FROM base AS builder
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY Makefile .
COPY .git ./.git

COPY go.* ./
COPY . ./
RUN make build


FROM ubuntu
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/gear-go /gear-go

ENTRYPOINT ["/gear-go", "serve"]