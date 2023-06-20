FROM golang:1.20

WORKDIR /src

COPY . .

RUN mkdir -p /src/data

RUN go get ./...

RUN go build -o /src/run ./cmd/main.go

CMD ["/src/run"]
