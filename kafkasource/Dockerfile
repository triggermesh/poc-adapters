FROM golang:1.17-buster AS builder

ENV CGO_ENABLED 1
ENV GOOS linux
ENV GOARCH amd64

ENV JENKINS_HOME=/var/jenkins_home

# Required to build in the kafka/confluent client
RUN mkdir $JENKINS_HOME && \
    chown ${uid}:${gid} $JENKINS_HOME

# RUN apt-get update && apt-get install  -y krb5-workstation krb5-libs krb5-auth-dialog libsas12-dev cryus-sasl-gssapi.x86_64 wget gcc glibc-static libstdc++-devel

RUN apt-get update && apt-get install -y librdkafka-dev  && \
    rm -rf /var/lib/apt/lists/*

RUN export HTTP_PRODY="http://vz-prozy.pncinit.net:8080" && \
    export HTTPS_PROXY="http://vz-prozy.pncint.net:8080" && \
    export http_proxy="http://vz-prozy.pncint.net:8080" && \
    export https_proxy="http://vz-prozy.pncint.net:8080" && \
    rm -rf /var/cache/* && \
    rm -rf /@System.solv && \
    groupadd --gid ${gid} ${group} && \
    useradd -u ${uid} -g ${group} ${user}

WORKDIR /project
COPY . ./
RUN cd /project/cmd && go build -o /project/bin/kafkatarget



FROM registry.access.redhat.com/ubi8/ubi-minimal
EXPOSE 8080
COPY --from=builder /project/bin/kafkatarget /kafkatarget

ENTRYPOINT ["/kafkatarget"]
