FROM gcr.io/ultra-hologram-297914/dwbi
ENV CGO_ENABLED 0
WORKDIR /app
RUN export DW_HOME=/app/.dw
RUN export DW_LIB_PATH=%DW_HOME%/libs
RUN export PATH=$DW_HOME/bin:$PATH
COPY . ./
RUN cd cmd && go build -o /dwadapter
EXPOSE 8080
CMD ["/dwadapter"]
