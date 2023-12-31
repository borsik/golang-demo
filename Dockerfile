FROM golang:1.20
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /golang-demo
EXPOSE 8080
CMD ["/golang-demo"]
