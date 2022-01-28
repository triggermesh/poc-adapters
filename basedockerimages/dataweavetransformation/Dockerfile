FROM ubuntu:latest
RUN mkdir /app
RUN apt-get update && apt-get install wget -y && apt-get install unzip
RUN mkdir /app/.dw
RUN cd /app/.dw && wget -c https://github.com/mulesoft-labs/data-weave-cli/releases/download/v1.0.16/dw-1.0.16-Linux && unzip dw-1.0.16-Linux
RUN chmod +x /app/.dw/bin/dw
RUN wget -c https://go.dev/dl/go1.17.6.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.17.6.linux-amd64.tar.gz
