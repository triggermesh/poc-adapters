FROM golang:1.17-alpine
RUN apk update
RUN apk upgrade
RUN apk add bash
# RUN set -ex && apk --no-cache add sudo
RUN mkdir /app
COPY . /app/dw
# RUN chmod +x /app/dw
RUN chmod g+rx,o+rx /app