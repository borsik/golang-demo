FROM golang:1.20
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /golang-demo
EXPOSE 8080
CMD ["/golang-demo"]