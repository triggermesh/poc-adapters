FROM golang:1.17-buster AS builder

ENV CGO_ENABLED 1
ENV GOOS linux
ENV GOARCH amd64



RUN apt-get update && apt-get install -y librdkafka-dev  && \
    rm -rf /var/lib/apt/lists/*


WORKDIR /project
COPY . ./
RUN cd /project/cmd && go build -o /project/bin/kafkatarget



FROM registry.access.redhat.com/ubi8/ubi-minimal
EXPOSE 8080
COPY --from=builder /project/bin/kafkatarget /kafkatarget

ENTRYPOINT ["/kafkatarget"]
