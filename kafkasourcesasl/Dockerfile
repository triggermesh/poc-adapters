FROM registry.access.redhat.com/ubi8/go-toolset:latest

USER root

COPY kafka.repo /etc/yum.repos.d/
RUN rpm --import https://packages.confluent.io/rpm/7.0/archive.key
RUN yum clean all
RUN go get -d -v github.com/confluentinc/confluent-kafka-go/kafka
RUN yum install  krb5-workstation pkg-config cyrus-sasl-gssapi.x86_64 librdkafka-devel -y

ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src

COPY go.mod .
COPY go.sum .
RUN go mod download

RUN git clone https://github.com/edenhill/librdkafka.git && cd librdkafka && ./configure --prefix /usr && make && make install && ldconfig
ENV PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/usr/lib/pkgconfig/

RUN go get -d -v github.com/confluentinc/confluent-kafka-go/kafka

COPY . .

RUN cd cmd && LDFLAGS='-extldflags=-w -s' go build -a -installsuffix cgo -tags dynamic  -o main


EXPOSE 8080

ENTRYPOINT ["/go/src/cmd/main"]
