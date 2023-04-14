FROM golang:1.16

WORKDIR /go/src/app
COPY . .
RUN go build -o bin/server cmd/server/main.go

CMD ["./bin/server"]