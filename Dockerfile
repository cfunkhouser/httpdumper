FROM golang:1
LABEL maintainer="christian@funkhouse.rs"

ENV GOBIN /bin

COPY . /httpdumper
WORKDIR /httpdumper

RUN go install

EXPOSE 8080
CMD [ "/bin/httpdumper" ]
