FROM golang:1.16

WORKDIR /go/src/app
COPY . .
RUN go build -o bin/server cmd/server/main.go

EXPOSE 5000
ENV PORT=5000

CMD ["./bin/server"]