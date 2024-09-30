# docker build . -t ghcr.io/reddio-com/reddio:latest
FROM golang:1.23-bookworm as builder

RUN mkdir /build
COPY . /build
RUN cd /build && git submodule init && git submodule update --recursive --checkout && make build

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && apt-get clean

# COPY ./conf /conf
RUN mkdir /reddio_db /yu
COPY --from=builder /build/reddio /reddio

CMD ["/reddio"]