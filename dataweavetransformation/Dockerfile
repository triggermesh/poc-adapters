FROM gcr.io/ultra-hologram-297914/dwbi:latest
WORKDIR /app
ENV CGO_ENABLED 0
ENV DW_HOME=/app/.dw
ENV DW_LIB_PATH=/app/.dw/libs
ENV PATH=/app/.dw/bin:$PATH
ENV PATH=$PATH:/usr/local/go/bin

COPY . ./
RUN cd cmd && go build -o /dwadapter
EXPOSE 8080
CMD ["/dwadapter"]
