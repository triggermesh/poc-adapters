FROM openjdk:18-jdk-alpine
COPY --from=golang:1.17-alpine /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"
RUN apk update && apk upgrade && apk add bash && apk add git
RUN mkdir /app
COPY dw /app/.dw/dw
RUN chmod +x /app/.dw/dw
RUN export DW_HOME=/app/.dw
RUN export DW_LIB_PATH=%DW_HOME%/libs
RUN export PATH=$DW_HOME/bin:$PATH
