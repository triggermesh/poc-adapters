FROM ubuntu:latest
RUN mkdir /app
COPY dw /app/.dw/bin/dw
RUN chmod +x /app/.dw/bin/dw
RUN apt-get update && apt-get install wget -y
RUN wget -c https://go.dev/dl/go1.17.6.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.17.6.linux-amd64.tar.gz
