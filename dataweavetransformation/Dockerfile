FROM gcr.io/ultra-hologram-297914/dwbi 
ENV CGO_ENABLED 0
WORKDIR /app
COPY . ./
RUN cd cmd && go build -o /dwadapter
EXPOSE 8080
CMD ["/dwadapter"]