FROM golang:1.17-buster AS builder
WORKDIR /project
COPY . ./
RUN cd /project/cmd && go build -o /project/bin/mongodbtarget

FROM registry.access.redhat.com/ubi8/ubi-minimal
EXPOSE 8080
COPY --from=builder /project/bin/mongodbtarget /mongodbtarget

ENTRYPOINT ["/mongodbtarget"]
