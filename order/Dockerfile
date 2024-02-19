FROM golang

WORKDIR /app

ENV GO111MODULE=on


COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./bin/main .

CMD ["/app/bin/main"]