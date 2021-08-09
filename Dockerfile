FROM golang:1.16-alpine
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY cybercafe.go ./
RUN go build -o /cybercafe
EXPOSE 8080
CMD ["/cybercafe"]