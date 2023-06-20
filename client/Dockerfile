FROM golang:1.20

WORKDIR /src

COPY . .

RUN go get ./...

RUN go build -o /src/run ./main.go

CMD ["/src/run"]
