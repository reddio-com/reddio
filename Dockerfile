# docker build . -t ghcr.io/reddio-com/reddio:latest
FROM golang:1.22-bookworm

RUN apt-get update && apt-get install -y ca-certificates && apt-get clean

COPY ./conf /conf
RUN mkdir /reddio_db /yu
COPY . /reddio

CMD ["/reddio"]
